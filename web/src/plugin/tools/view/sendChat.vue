<template>
  <div>
    <div class="gva-table-box">
      <el-card class="box-card mb-4">
        <el-form :inline="true" :model="createForm" label-width="120px">
          <el-form-item label="房间ID">
            <el-input v-model="createForm.roomId" placeholder="请输入房间ID" style="width:200px" clearable />
          </el-form-item>
          <el-form-item label="选择环境">
            <el-select v-model="createForm.environmentKey" placeholder="选择环境" filterable style="width:220px">
              <el-option v-for="env in environmentOptions" :key="env.key" :label="env.name" :value="env.key" />
            </el-select>
          </el-form-item>
          <el-form-item label="账号数量">
            <el-input-number v-model="createForm.accountCount" :min="1" :max="100" style="width:150px" />
          </el-form-item>
          <el-form-item label="每账号消息数">
            <el-input-number v-model="createForm.msgCountPerAccount" :min="1" :max="10000" style="width:150px" />
          </el-form-item>
          <el-form-item label="消息间隔(毫秒)">
            <el-input-number v-model="createForm.msgInterval" :min="0" :max="60000" style="width:150px" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="creating" @click="handleCreate">创建任务</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <div class="gva-table-box">
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
          <el-table-column align="left" label="房间ID" prop="roomId" width="120" />
          <el-table-column align="left" label="环境Key" prop="environmentKey" width="180" show-overflow-tooltip />
          <el-table-column align="left" label="账号数" prop="accountCount" width="80" />
          <el-table-column align="left" label="消息数/账号" prop="msgCountPerAccount" width="110" />
          <el-table-column align="left" label="间隔(毫秒)" prop="msgInterval" width="100" />
          <el-table-column align="left" label="状态" prop="status" width="100">
            <template #default="scope">
              <el-tag :type="statusType(scope.row.status)">{{ statusLabel(scope.row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column align="left" label="成功数量" prop="successCount" width="100" />
          <el-table-column align="left" label="操作" fixed="right" width="120">
            <template #default="scope">
              <el-button
                v-if="scope.row.status === 'running'"
                type="danger"
                link
                icon="close"
                @click="handleStop(scope.row)"
              >停止</el-button>
              <span v-else class="text-gray-400">-</span>
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
    </div>
  </div>
</template>

<script setup>
import { createSendChatTask, getSendChatTaskList, stopSendChatTask } from '@/plugin/tools/api/sendChat'
import { getEnvironmentList } from '@/plugin/tools/api/environment'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, onMounted, onUnmounted } from 'vue'

defineOptions({ name: 'ToolsSendChat' })

const createForm = ref({
  roomId: '',
  environmentKey: '',
  accountCount: 10,
  msgCountPerAccount: 100,
  msgInterval: 100
})

const environmentOptions = ref([])
const creating = ref(false)
const stoppingIds = ref(new Set())

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

let pollTimer = null

const statusType = (status) => {
  const map = { running: 'warning', completed: 'success', stopped: 'info', failed: 'danger' }
  return map[status] || ''
}

const statusLabel = (status) => {
  const map = { running: '运行中', completed: '已完成', stopped: '已停止', failed: '失败' }
  return map[status] || status
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const getTableData = async () => {
  const res = await getSendChatTaskList({
    page: page.value,
    pageSize: pageSize.value
  })
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

const startPolling = () => {
  stopPolling()
  pollTimer = setInterval(() => {
    if (tableData.value.some(r => r.status === 'running')) {
      getTableData()
    }
  }, 3000)
}

const stopPolling = () => {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

onMounted(() => {
  getTableData()
  loadEnvironments()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})

const handleCreate = async () => {
  if (!createForm.value.roomId) {
    ElMessage.warning('请输入房间ID')
    return
  }
  if (!createForm.value.environmentKey) {
    ElMessage.warning('请选择环境')
    return
  }
  creating.value = true
  try {
    const res = await createSendChatTask(createForm.value)
    if (res.code === 0) {
      ElMessage.success('任务已创建')
      page.value = 1
      getTableData()
    } else {
      ElMessage.error(res.msg || '创建失败')
    }
  } finally {
    creating.value = false
  }
}

const handleStop = async (row) => {
  ElMessageBox.confirm('确定要停止该任务吗?', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    .then(async () => {
      stoppingIds.value.add(row.ID)
      try {
        const res = await stopSendChatTask({ ID: row.ID })
        if (res.code === 0) {
          ElMessage.success('任务已停止')
          getTableData()
        } else {
          ElMessage.error(res.msg || '停止失败')
        }
      } finally {
        stoppingIds.value.delete(row.ID)
      }
    })
}
</script>
