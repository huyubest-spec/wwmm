package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	loginURL := "http://localhost:8080/api/user/login"
	uploadURL := "http://localhost:8080/api/photo/upload"

	token := login(loginURL, "photographer", "photo123")
	if token == "" {
		token = login(loginURL, "alice", "alice123")
	}
	if token == "" {
		log.Fatal("摄影师登录失败")
	}

	photos := []struct {
		Title       string
		Description string
		Category    string
		Location    string
		ShootTime   string
		CameraInfo  string
		ColorR      uint8
		ColorG      uint8
		ColorB      uint8
	}{
		{"云海日出 · 衡山", "凌晨四点的南岳衡山，金色的阳光穿透云海，恍若仙境。", "风光", "湖南·衡山", "2025-10-12", "Sony A7M4 + 24-70GM F2.8", 255, 200, 100},
		{"湘江夜色", "橘子洲头的霓虹灯倒映在江面上，光与影的交织。", "城市", "湖南·长沙", "2025-09-30", "Canon R6 + 50mm F1.4", 60, 120, 200},
		{"古镇小巷", "凤凰古城的青石板路，记录老街的人间烟火。", "纪实", "湖南·凤凰", "2025-08-15", "Fuji X-T5 + 35mm F1.4", 180, 120, 80},
		{"田野花开", "浏阳大围山的油菜花海，春风十里。", "风光", "湖南·浏阳", "2025-04-08", "Nikon Z7 + 70-200 F2.8", 250, 220, 80},
		{"星轨 · 草原之夜", "内蒙古乌兰布统，相机长曝光下北半球的星轨。", "风光", "内蒙古·乌兰布统", "2025-07-22", "Sony A7S3 + 14mm F1.8", 30, 50, 100},
		{"人间烟火", "菜市场清晨的烟火气，最抚凡人心。", "纪实", "湖南·长沙", "2025-06-10", "Leica Q2", 220, 150, 90},
	}

	tmpDir := "./test_photos"
	_ = os.RemoveAll(tmpDir)
	_ = os.MkdirAll(tmpDir, 0755)

	for i, p := range photos {
		path := filepath.Join(tmpDir, fmt.Sprintf("photo_%d.png", i+1))
		if err := writeColoredPNG(path, 800, 600, p.ColorR, p.ColorG, p.ColorB, p.Title); err != nil {
			log.Printf("生成图片失败: %v", err)
			continue
		}
		if err := upload(uploadURL, token, path, p.Title, p.Description, p.Category, p.Location, p.ShootTime, p.CameraInfo); err != nil {
			log.Printf("上传 [%s] 失败: %v", p.Title, err)
		} else {
			log.Printf("OK [%s] 已上传", p.Title)
		}
	}

	// 模拟一些用户投票，让排行榜有内容
	voters := []struct{ name, pwd string }{
		{"voter", "vote123"},
		{"admin", "admin123"},
	}
	for _, v := range voters {
		tok := login(loginURL, v.name, v.pwd)
		if tok == "" {
			continue
		}
		for pid := 1; pid <= 6; pid++ {
			url := fmt.Sprintf("http://localhost:8080/api/photo/%d/vote", pid)
			req, _ := http.NewRequest("POST", url, nil)
			req.Header.Set("Token", tok)
			req.Header.Set("Authorization", "Bearer "+tok)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				continue
			}
			resp.Body.Close()
		}
	}

	log.Println("全部完成")
}

func login(url, u, p string) string {
	body := fmt.Sprintf(`{"username":%q,"password":%q}`, u, p)
	resp, err := http.Post(url, "application/json", bytes.NewBufferString(body))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ""
	}
	var r struct {
		Code int
		Data struct {
			Token string
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return ""
	}
	if r.Code != 0 {
		return ""
	}
	return r.Data.Token
}

func upload(url, token, filePath, title, desc, cat, loc, shoot, cam string) error {
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	mw.WriteField("title", title)
	mw.WriteField("description", desc)
	mw.WriteField("category", cat)
	mw.WriteField("shootLocation", loc)
	mw.WriteField("shootTime", shoot)
	mw.WriteField("cameraInfo", cam)
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	fw, _ := mw.CreateFormFile("image", filepath.Base(filePath))
	if _, err := io.Copy(fw, f); err != nil {
		return err
	}
	mw.Close()

	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Token", token)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("status %d, body: %s", resp.StatusCode, string(raw))
	}
	log.Printf("  → %s", string(raw)[:min(120, len(raw))])
	return nil
}

func min(a, b int) int { if a < b { return a }; return b }

func writeColoredPNG(path string, w, h int, r, g, b uint8, label string) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			fx := float64(x) / float64(w)
			fy := float64(y) / float64(h)
			cr := uint8(float64(r) * (1 - fx*0.3) * (1 - fy*0.3))
			cg := uint8(float64(g) * (1 - fx*0.3) * (1 - fy*0.3))
			cb := uint8(float64(b) * (1 - fx*0.3) * (1 - fy*0.3))
			img.SetRGBA(x, y, color.RGBA{cr, cg, cb, 255})
		}
	}
	for i := 0; i < w+h; i += 12 {
		drawLine(img, 0, i, i, 0, color.RGBA{255, 255, 255, 50})
	}
	cx, cy, radius := w/2, h/2, 120
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				img.SetRGBA(cx+x, cy+y, color.RGBA{255, 255, 255, 180})
			}
		}
	}
	for x := 0; x < w; x++ {
		img.SetRGBA(x, 0, color.RGBA{255, 255, 255, 255})
		img.SetRGBA(x, h-1, color.RGBA{255, 255, 255, 255})
	}
	for y := 0; y < h; y++ {
		img.SetRGBA(0, y, color.RGBA{255, 255, 255, 255})
		img.SetRGBA(w-1, y, color.RGBA{255, 255, 255, 255})
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.RGBA) {
	dx := abs2(x2 - x1)
	dy := abs2(y2 - y1)
	sx, sy := 1, 1
	if x1 >= x2 {
		sx = -1
	}
	if y1 >= y2 {
		sy = -1
	}
	err0 := dx - dy
	x, y := x1, y1
	for {
		if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
			img.SetRGBA(x, y, c)
		}
		if x == x2 && y == y2 {
			break
		}
		e2 := 2 * err0
		if e2 > -dy {
			err0 -= dy
			x += sx
		}
		if e2 < dx {
			err0 += dx
			y += sy
		}
	}
}

func abs2(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
