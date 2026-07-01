<template>
  <div style="max-width: 720px; margin: 0 auto;">
    <el-page-header @back="$router.push('/')" content="上传作品" style="margin-bottom: 16px;"></el-page-header>
    <el-card>
      <el-alert type="info" :closable="false" show-icon style="margin-bottom: 16px;">
        <div slot="title">⛓ 区块链存证提示</div>
        上传时系统会自动计算图片的 SHA-256 哈希，连同作者、时间、描述等元数据写入区块链。
        一旦上链便不可篡改，是您作品归属的强有力证明。
      </el-alert>

      <el-form :model="form" :rules="rules" ref="form" label-width="100px" size="medium">
        <el-form-item label="作品标题" prop="title">
          <el-input v-model="form.title" placeholder="给你的作品起一个响亮的名字"></el-input>
        </el-form-item>
        <el-form-item label="作品描述" prop="description">
          <el-input type="textarea" :rows="4" v-model="form.description" placeholder="创作背景、灵感、构图思路等"></el-input>
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="分类">
              <el-select v-model="form.category" placeholder="请选择" style="width: 100%;">
                <el-option label="风光摄影" value="风光"></el-option>
                <el-option label="人像摄影" value="人像"></el-option>
                <el-option label="纪实摄影" value="纪实"></el-option>
                <el-option label="街拍摄影" value="街拍"></el-option>
                <el-option label="微距摄影" value="微距"></el-option>
                <el-option label="其他" value="其他"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="拍摄地点">
              <el-input v-model="form.shootLocation" placeholder="选填"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="拍摄时间">
              <el-date-picker v-model="form.shootTime" type="date" value-format="yyyy-MM-dd" placeholder="选填" style="width: 100%;"></el-date-picker>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="器材信息">
              <el-input v-model="form.cameraInfo" placeholder="如 Sony A7M4 + 24-70GM"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="作品图片" prop="file">
          <el-upload
            ref="upload"
            :auto-upload="false"
            :limit="1"
            :on-change="onFileChange"
            :on-exceed="onExceed"
            drag
            action="">
            <i class="el-icon-upload"></i>
            <div class="el-upload__text">将图片拖到此处，或<em>点击上传</em></div>
            <div slot="tip" class="el-upload__tip">支持 jpg / png 格式，大小不超过 10MB</div>
          </el-upload>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="onSubmit" style="width: 100%;">
            <i class="el-icon-upload"></i>
            提交存证（写入区块链）
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'UploadPhoto',
  data() {
    return {
      loading: false,
      form: {
        title: '',
        description: '',
        category: '',
        shootLocation: '',
        shootTime: '',
        cameraInfo: '',
        file: null
      },
      rules: {
        title: [{ required: true, message: '请输入作品标题', trigger: 'blur' }],
        file: [{ required: true, message: '请选择图片' }]
      }
    }
  },
  methods: {
    onFileChange(file) {
      this.form.file = file.raw
    },
    onExceed() {
      this.$message.warning('只允许上传 1 张图片')
    },
    async onSubmit() {
      this.$refs.form.validate(async valid => {
        if (!valid) return
        const fd = new FormData()
        fd.append('title', this.form.title)
        fd.append('description', this.form.description)
        fd.append('category', this.form.category)
        fd.append('shootLocation', this.form.shootLocation)
        fd.append('shootTime', this.form.shootTime)
        fd.append('cameraInfo', this.form.cameraInfo)
        fd.append('image', this.form.file)
        this.loading = true
        try {
          const r = await this.$api.post('/photo/upload', fd, { headers: { 'Content-Type': 'multipart/form-data' } })
          this.$message.success('提交成功！存证信息已写入区块链')
          this.$alert('存证交易哈希：' + (r.data.chainTxHash || '已生成'), '上链成功', {
            confirmButtonText: '查看作品',
            callback: () => this.$router.push('/photo/' + r.data.photoId)
          })
        } catch (e) {
          this.$message.error(e.message || '提交失败')
        } finally {
          this.loading = false
        }
      })
    }
  }
}
</script>
