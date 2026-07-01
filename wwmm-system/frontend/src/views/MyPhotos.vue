<template>
  <div v-loading="loading">
    <el-page-header @back="$router.push('/')" content="我的作品" style="margin-bottom: 16px;"></el-page-header>
    <h2 class="section-title">我上传的作品</h2>
    <div v-if="!photos.length" class="empty-hint">还没有上传过作品</div>
    <el-table v-else :data="photos" stripe>
      <el-table-column label="封面" width="120">
        <template slot-scope="s">
          <div :style="{ width:'80px', height:'60px', backgroundImage:'url(' + getImg(s.row.imageUrl) + ')', backgroundSize:'cover', backgroundPosition:'center', borderRadius:'4px' }"></div>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题" min-width="160"></el-table-column>
      <el-table-column label="状态" width="100">
        <template slot-scope="s">
          <el-tag :type="statusTag(s.row.status)" size="mini" effect="plain">{{ statusLabel(s.row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="上链" width="100">
        <template slot-scope="s">
          <el-tag v-if="s.row.isOnChain" type="success" size="mini">⛓ 已上链</el-tag>
          <el-tag v-else type="info" size="mini">未上链</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="voteCount" label="得票" width="80"></el-table-column>
      <el-table-column prop="viewCount" label="浏览" width="80"></el-table-column>
      <el-table-column prop="createdAt" label="创建时间" width="170"></el-table-column>
      <el-table-column label="操作" width="180">
        <template slot-scope="s">
          <el-button size="mini" @click="$router.push('/photo/' + s.row.photoId)">查看</el-button>
          <el-button size="mini" type="text" v-if="s.row.chainTxHash" @click="$router.push('/chain/tx/' + s.row.chainTxHash)">查看链交易</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
export default {
  name: 'MyPhotos',
  data() {
    return { loading: true, photos: [] }
  },
  async mounted() { await this.load() },
  methods: {
    statusLabel(s) { return ['待审核', '已通过', '已拒绝'][s] || '-' },
    statusTag(s) { return ['info', 'success', 'danger'][s] || 'info' },
    getImg(url) {
      if (!url) return ''
      if (url.startsWith('http')) return url
      const base = (this.$api.defaults.baseURL || '').replace('/api', '')
      return base + url
    },
    async load() {
      this.loading = true
      try {
        const r = await this.$api.get('/photo/mine')
        this.photos = r.data.list || []
      } catch (e) {
        this.$message.error(e.message)
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
