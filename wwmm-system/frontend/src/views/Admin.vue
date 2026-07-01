<template>
  <div v-loading="loading">
    <el-page-header @back="$router.push('/')" content="审核中心" style="margin-bottom: 16px;"></el-page-header>
    <h2 class="section-title">待审核作品</h2>
    <el-alert type="warning" :closable="false" show-icon style="margin-bottom: 16px;">
      摄影师提交的所有作品均需要经过审核才能展示给公众。
      审核通过后，作品才会被允许进入投票环节。
    </el-alert>
    <div v-if="!photos.length" class="empty-hint">没有待审核的作品</div>
    <el-table v-else :data="photos" stripe>
      <el-table-column label="封面" width="120">
        <template slot-scope="s">
          <div :style="{ width:'80px', height:'60px', backgroundImage:'url(' + getImg(s.row.imageUrl) + ')', backgroundSize:'cover', backgroundPosition:'center', borderRadius:'4px' }"></div>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题" min-width="160"></el-table-column>
      <el-table-column label="作者" width="140">
        <template slot-scope="s">{{ s.row.photographerName }}</template>
      </el-table-column>
      <el-table-column label="图片哈希" min-width="200">
        <template slot-scope="s">
          <span class="hash-cell">{{ (s.row.imageHash || '').substr(0, 24) }}…</span>
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="提交时间" width="170"></el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template slot-scope="s">
          <el-button size="mini" @click="$router.push('/photo/' + s.row.photoId)">预览</el-button>
          <el-button size="mini" type="success" @click="audit(s.row.photoId, true)">通过</el-button>
          <el-button size="mini" type="danger" @click="audit(s.row.photoId, false)">拒绝</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
export default {
  name: 'Admin',
  data() { return { loading: true, photos: [] } },
  async mounted() { await this.load() },
  methods: {
    getImg(url) {
      if (!url) return ''
      if (url.startsWith('http')) return url
      const base = (this.$api.defaults.baseURL || '').replace('/api', '')
      return base + url
    },
    async load() {
      this.loading = true
      try {
        const r = await this.$api.get('/photo/pending')
        this.photos = r.data.list || []
      } catch (e) {
        this.$message.error(e.message)
      } finally { this.loading = false }
    },
    audit(id, approve) {
      this.$prompt('请填写审核意见', '审核', {
        confirmButtonText: approve ? '通过' : '拒绝',
        cancelButtonText: '取消',
        inputPattern: /.+/,
        inputErrorMessage: '请填写意见'
      }).then(({ value }) => {
        this.$api.post('/photo/' + id + '/audit', { approve, comment: value }).then(() => {
          this.$message.success('审核完成')
          this.load()
        }).catch(err => this.$message.error(err.message))
      }).catch(() => {})
    }
  }
}
</script>
