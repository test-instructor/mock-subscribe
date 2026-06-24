<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="环境" prop="environmentKey">
          <el-select v-model="searchInfo.environmentKey" placeholder="选择环境" clearable filterable style="width:200px">
            <el-option v-for="env in environmentOptions" :key="env.key" :label="env.name" :value="env.key" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户ID" prop="userId">
          <el-input v-model.number="searchInfo.userId" placeholder="单个用户ID" clearable style="width:160px" />
        </el-form-item>
        <el-form-item label="批量ID" prop="userIdsText">
          <el-input
            v-model="searchInfo.userIdsText"
            type="textarea"
            :rows="2"
            placeholder="批量查询(每行一个用户ID)"
            style="width:280px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
      >
        <el-table-column align="left" label="日期" prop="CreatedAt" width="180">
          <template #default="scope">
            <span>{{ formatDate(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="环境Key" prop="environmentKey" width="200" show-overflow-tooltip />
        <el-table-column align="left" label="用户ID" prop="userId" width="150" />
        <el-table-column align="left" label="操作" fixed="right" width="120">
          <template #default="scope">
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    <el-drawer
      v-model="dialogFormVisible"
      destroy-on-close
      size="600px"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">批量新增用户</span>
          <div>
            <el-button type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>
      <el-form ref="elFormRef" :model="formData" label-position="top" :rules="rule" label-width="80px">
        <el-form-item label="环境" prop="environmentKey">
          <el-select v-model="formData.environmentKey" placeholder="选择环境" filterable style="width:100%">
            <el-option v-for="env in environmentOptions" :key="env.key" :label="env.name + ' (' + env.key + ')'" :value="env.key" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户ID列表" prop="userIds">
          <el-input
            v-model="formData.userIds"
            type="textarea"
            :rows="10"
            placeholder="一行一个用户ID"
          />
          <div class="text-sm text-gray-500 mt-1">每行输入一个用户ID</div>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { getUserRelationList, createUserRelation, deleteUserRelation, getUserIdsByEnvironment } from '@/plugin/tools/api/userRelation'
import { getEnvironmentList } from '@/plugin/tools/api/environment'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'

defineOptions({ name: 'ToolsUserRelation' })

const formData = ref({
  environmentKey: '',
  userIds: ''
})

const rule = {
  environmentKey: [{ required: true, message: '请选择环境', trigger: 'change' }],
  userIds: [{ required: true, message: '请输入用户ID', trigger: 'blur' }],
}

const searchInfo = ref({
  environmentKey: '',
  userId: '',
  userIdsText: ''
})

const environmentOptions = ref([])

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

const parseBatchIds = (text) => {
  if (!text) return []
  const out = []
  const lines = String(text).split(/[\n,; \t]+/)
  for (const line of lines) {
    const v = line.trim()
    if (!v) continue
    const n = Number(v)
    if (Number.isFinite(n) && n > 0) out.push(n)
  }
  return out
}

const onReset = () => {
  searchInfo.value = { environmentKey: '', userId: '', userIdsText: '' }
  page.value = 1
  getTableData()
}

const onSubmit = () => {
  page.value = 1
  getTableData()
}

const elSearchFormRef = ref(null)
const elFormRef = ref(null)

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const getTableData = async () => {
  const userIds = parseBatchIds(searchInfo.value.userIdsText)
  const params = {
    page: page.value,
    pageSize: pageSize.value,
    environmentKey: searchInfo.value.environmentKey || undefined,
    userId: searchInfo.value.userId ? Number(searchInfo.value.userId) : undefined
  }
  if (userIds.length > 0) {
    params['userIds[]'] = userIds
  }
  const res = await getUserRelationList(params)
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}

const loadEnvironments = async () => {
  const res = await getEnvironmentList({ page: 1, pageSize: 1000 })
  if (res.code === 0) {
    environmentOptions.value = res.data.list
  }
}

onMounted(() => {
  getTableData()
  loadEnvironments()
})

const dialogFormVisible = ref(false)

const openDialog = () => {
  dialogFormVisible.value = true
}

const closeDialog = () => {
  dialogFormVisible.value = false
  elFormRef.value?.resetFields()
  formData.value = { environmentKey: '', userIds: '' }
}

const deleteRow = async (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    .then(async () => {
      const res = await deleteUserRelation({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '删除成功' })
        getTableData()
      }
    })
}

const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (valid) {
      const res = await createUserRelation(formData.value)
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '创建成功' })
        closeDialog()
        getTableData()
      }
    }
  })
}
</script>