<template>
  <div v-loading="loading">
    <el-page-header @back="$router.push('/')" content="作品详情" style="margin-bottom: 16px;"></el-page-header>
    <el-row :gutter="24" v-if="photo.photoId">
      <el-col :span="14">
        <el-card>
          <div class="photo-img" :style="{ backgroundImage: 'url(' + imageFullUrl + ')'}"></div>
        </el-card>
        <el-card style="margin-top: 16px;">
          <div slot="header">
            <span style="font-weight: 600;">🔐 区块链存证信息</span>
            <el-tag v-if="photo.isOnChain" type="success" size="mini" effect="dark" style="margin-left: 8px;">已上链</el-tag>
            <el-tag v-else type="info" size="mini" effect="plain" style="margin-left: 8px;">未上链</el-tag>
          </div>
          <el-descriptions :column="1" size="small" border>
            <el-descriptions-item label="图片 SHA-256">
              <span class="hash-cell">{{ photo.imageHash }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="上链交易哈希">
              <span class="hash-cell" v-if="photo.chainTxHash">
                <router-link :to="'/chain/tx/' + photo.chainTxHash" style="color:#2c5fe0;">{{ photo.chainTxHash }}</router-link>
              </span>
              <span v-else>-</span>
            </el-descriptions-item>
            <el-descriptions-item label="文件大小">{{ (photo.fileSize/1024).toFixed(2) }} KB</el-descriptions-item>
            <el-descriptions-item label="提交时间">{{ photo.createdAt }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      <el-col :span="10">
        <el-card>
          <h2 style="margin: 0 0 4px 0; font-size: 22px; font-weight: 700;">{{ photo.title }}</h2>
          <div style="color: #8a96b3; font-size: 13px; margin-bottom: 14px;">
            📷 {{ photo.photographerName }}
            <el-tag size="mini" :type="photo.status===1?'success':'info'" effect="plain" style="margin-left: 6px;">
              {{ statusLabel(photo.status) }}
            </el-tag>
          </div>
          <div style="line-height: 1.8; color: #4a5568; font-size: 14px; margin-bottom: 16px;">{{ photo.description || '该作品暂无描述' }}</div>

          <el-descriptions :column="2" size="small" style="margin-bottom: 16px;">
            <el-descriptions-item label="分类">{{ photo.category || '-' }}</el-descriptions-item>
            <el-descriptions-item label="拍摄地">{{ photo.shootLocation || '-' }}</el-descriptions-item>
            <el-descriptions-item label="拍摄时间">{{ photo.shootTime || '-' }}</el-descriptions-item>
            <el-descriptions-item label="器材">{{ photo.cameraInfo || '-' }}</el-descriptions-item>
          </el-descriptions>

          <div class="vote-bar">
            <div>
              <div class="vote-num">❤ {{ photo.voteCount || 0 }}</div>
              <div class="vote-label">累计票数</div>
            </div>
            <div>
              <div class="vote-num">👁 {{ photo.viewCount || 0 }}</div>
              <div class="vote-label">浏览数</div>
            </div>
          </div>

          <el-button v-if="!isOwner && photo.status===1" :type="photo.hasVoted?'info':'primary'" size="large" style="width: 100%; margin-top: 8px;" :loading="voting" :disabled="photo.hasVoted" @click="onVote">
            <i class="el-icon-thumb"></i>
            {{ photo.hasVoted ? '您已为此作品投过票' : '为它投一票（数据将上链）' }}
          </el-button>
          <el-alert v-if="!user.userId" type="warning" :closable="false" show-icon style="margin-top: 8px;"
            title="请先登录后再投票" />
          <el-alert v-if="isOwner" type="info" :closable="false" show-icon style="margin-top: 8px;"
            title="这是您自己的作品，无法给自己投票" />
          <el-alert v-if="voteTxHash" type="success" :closable="false" show-icon style="margin-top: 12px;">
            <div slot="title">投票成功！交易已上链</div>
            <div style="margin-top: 4px; font-size: 12px;">交易哈希：
              <router-link :to="'/chain/tx/' + voteTxHash" class="hash-cell" style="color:#2c5fe0;">{{ voteTxHash }}</router-link>
            </div>
          </el-alert>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: 'PhotoDetail',
  props: ['id'],
  data() {
    return {
      loading: true,
      voting: false,
      photo: {},
      user: {},
      voteTxHash: ''
    }
  },
  computed: {
    isOwner() {
      return this.user.userId && this.user.userId === this.photographerId
    },
    photographerId() {
      return this.photo.photographerId
    },
    imageFullUrl() {
      if (!this.photo.imageUrl) return ''
      if (this.photo.imageUrl.startsWith('http')) return this.photo.imageUrl
      const base = (this.$api.defaults.baseURL || '').replace('/api', '')
      return base + this.photo.imageUrl
    }
  },
  async mounted() {
    this.loadUser()
    await this.load()
  },
  methods: {
    loadUser() {
      const raw = localStorage.getItem('wwmm_user')
      if (raw) try { this.user = JSON.parse(raw) } catch (e) {}
    },
    statusLabel(s) {
      return ['待审核', '已通过', '已拒绝'][s] || '-'
    },
    async load() {
      this.loading = true
      try {
        const r = await this.$api.get('/photo/' + this.id)
        this.photo = r.data
      } catch (e) {
        this.$message.error('加载失败：' + e.message)
      } finally {
        this.loading = false
      }
    },
    async onVote() {
      if (!this.user.userId) { this.$router.push('/login'); return }
      this.voting = true
      try {
        const r = await this.$api.post('/photo/' + this.id + '/vote')
        this.voteTxHash = r.data.txHash
        this.$message.success('投票成功，已写入区块链')
        await this.load()
      } catch (e) {
        this.$message.error(e.message || '投票失败')
      } finally {
        this.voting = false
      }
    }
  }
}
</script>

<style scoped>
.photo-img {
  width: 100%;
  height: 480px;
  background-color: #f0f2f7;
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
  border-radius: 6px;
}
.vote-bar { display: flex; gap: 16px; padding: 14px 0; border-top: 1px solid #e6ebf5; border-bottom: 1px solid #e6ebf5; margin-bottom: 12px; }
.vote-bar > div { flex: 1; text-align: center; }
.vote-num { font-size: 22px; font-weight: 700; color: #d44030; }
.vote-label { font-size: 12px; color: #8a96b3; margin-top: 2px; }
</style>
