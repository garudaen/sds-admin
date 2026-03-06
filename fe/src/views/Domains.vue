<template>
  <div class="domains-page">
    <el-tabs v-model="activeTab" class="tabs">
      <el-tab-pane label="全部域名" name="all">
        <template #label>
          <span>全部域名</span>
          <el-badge :value="domains.length" class="tab-badge" />
        </template>
      </el-tab-pane>
      <el-tab-pane label="已启用" name="active" />
      <el-tab-pane label="已禁用" name="disabled" />
    </el-tabs>

    <div class="toolbar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索域名..."
        clearable
        class="search-input"
        :prefix-icon="Search"
      />
      <el-select v-model="statusFilter" placeholder="全部状态" clearable class="filter-select">
        <el-option label="已启用" value="active" />
        <el-option label="已禁用" value="disabled" />
      </el-select>
    </div>

    <div class="table-wrap">
      <el-table
        v-loading="loading"
        :data="filteredDomains"
        stripe
        style="width: 100%"
        :empty-text="'暂无域名数据，点击右上角「添加域名」创建'"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="域名" min-width="160">
          <template #default="{ row }">
            <span class="domain-name">{{ row.domainName }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="140" show-overflow-tooltip>
          <template #default="{ row }">{{ row.description || '-' }}</template>
        </el-table-column>
        <el-table-column label="递归" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.recursive ? 'success' : 'info'" size="small">
              {{ row.recursive ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '已启用' : '已禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="goToRecords(row)">解析</el-button>
            <el-button type="primary" link size="small" @click="openEditModal(row)">编辑</el-button>
            <template v-if="row.status === 'active'">
              <el-button type="warning" link size="small" @click="disableDomain(row.id)">禁用</el-button>
            </template>
            <template v-else>
              <el-button type="success" link size="small" @click="enableDomain(row.id)">启用</el-button>
            </template>
            <el-button
              type="danger"
              link
              size="small"
              :loading="deletingDomain === row.id"
              @click="deleteDomain(row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div v-if="filteredDomains.length > 0" class="pagination">
      <span class="total">共 {{ filteredDomains.length }} 条</span>
    </div>

    <!-- 编辑域名 -->
    <el-dialog
      v-model="showEditModal"
      title="编辑域名配置"
      width="480px"
      destroy-on-close
      @close="closeEditModal"
    >
      <el-form :model="editFormData" label-width="120px">
        <el-form-item label="域名名称">
          <el-input :model-value="editingDomain?.domainName" disabled />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="editFormData.description" placeholder="请输入域名描述" clearable />
        </el-form-item>
        <el-form-item label="递归查询">
          <el-switch v-model="editFormData.recursive" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeEditModal">取消</el-button>
        <el-button type="primary" :loading="updatingConfig" @click="updateDomainConfig">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../api/request'

const router = useRouter()
const store = useStore()

const activeTab = ref('all')
const searchKeyword = ref('')
const statusFilter = ref('')
const domains = ref([])
const loading = ref(false)
const deletingDomain = ref(null)
const showEditModal = ref(false)
const editingDomain = ref(null)
const editFormData = ref({ description: '', recursive: false })
const updatingConfig = ref(false)

const filteredDomains = computed(() => {
  let list = domains.value
  if (activeTab.value === 'active') list = list.filter((d) => d.status === 'active')
  else if (activeTab.value === 'disabled') list = list.filter((d) => d.status === 'disabled')
  if (statusFilter.value) list = list.filter((d) => d.status === statusFilter.value)
  if (searchKeyword.value) {
    const k = searchKeyword.value.toLowerCase()
    list = list.filter(
      (d) =>
        d.domainName.toLowerCase().includes(k) ||
        (d.description && d.description.toLowerCase().includes(k))
    )
  }
  return list
})

async function fetchDomains() {
  loading.value = true
  try {
    const res = await request.get('/domains')
    if (res.data?.code === 200) {
      domains.value = res.data.data || []
      store.dispatch('domains/setList', domains.value)
    }
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function goToRecords(domain) {
  store.dispatch('app/goToRecords', domain)
  router.push(`/domains/${domain.id}/records`)
}

function openEditModal(domain) {
  editingDomain.value = domain
  editFormData.value = {
    description: domain.description || '',
    recursive: domain.recursive,
  }
  showEditModal.value = true
}

function closeEditModal() {
  showEditModal.value = false
  editingDomain.value = null
}

async function updateDomainConfig() {
  if (!editingDomain.value) return
  updatingConfig.value = true
  try {
    await request.patch(`/domains/${editingDomain.value.id}/config`, editFormData.value)
    await fetchDomains()
    closeEditModal()
    ElMessage.success('域名配置更新成功')
  } catch (_) {
    // 错误已在 request 拦截器处理
  } finally {
    updatingConfig.value = false
  }
}

async function deleteDomain(id) {
  await ElMessageBox.confirm('确定要删除该域名吗？', '提示', {
    type: 'warning',
  }).catch(() => {})
  deletingDomain.value = id
  try {
    await request.delete(`/domains/${id}`)
    await fetchDomains()
    ElMessage.success('删除成功')
  } finally {
    deletingDomain.value = null
  }
}

async function disableDomain(id) {
  await request.post(`/domains/${id}/disable`)
  await fetchDomains()
  ElMessage.success('已禁用')
}

async function enableDomain(id) {
  await request.post(`/domains/${id}/enable`)
  await fetchDomains()
  ElMessage.success('已启用')
}

onMounted(() => {
  fetchDomains()
  window.addEventListener('domains-refresh', fetchDomains)
})

onUnmounted(() => {
  window.removeEventListener('domains-refresh', fetchDomains)
})
</script>

<style scoped>
.domains-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.tabs {
  padding: 0 20px;
}

.tab-badge {
  margin-left: 6px;
}

.toolbar {
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
}

.search-input {
  max-width: 280px;
}

.filter-select {
  width: 120px;
}

.table-wrap {
  flex: 1;
  overflow: auto;
  padding: 0 20px;
}

.domain-name {
  font-weight: 500;
  color: #303133;
}

.pagination {
  padding: 12px 20px;
  border-top: 1px solid #e4e7ed;
}

.total {
  color: #606266;
  font-size: 14px;
}

@media (max-width: 768px) {
  .toolbar {
    flex-direction: column;
    padding: 12px;
  }

  .search-input {
    max-width: none;
  }

  .filter-select {
    width: 100%;
  }

  .table-wrap {
    padding: 12px;
  }
}
</style>
