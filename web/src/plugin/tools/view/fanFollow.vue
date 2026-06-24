<template>
  <div>
    <div class="gva-table-box">
      <el-card class="box-card mb-4">
        <el-form :inline="true" :model="formData" label-width="100px">
          <el-form-item label="选择环境">
            <el-select v-model="formData.environmentKey" placeholder="选择环境" filterable style="width:220px" @change="onEnvChange">
              <el-option v-for="env in environmentOptions" :key="env.key" :label="env.name" :value="env.key" />
            </el-select>
          </el-form-item>
          <el-form-item label="用户ID">
            <el-input-number v-model="formData.userId" :min="1" placeholder="目标用户ID" style="width:200px" />
          </el-form-item>
          <el-form-item label="操作类型">
            <el-radio-group v-model="formData.operation">
              <el-radio label="fans">粉丝</el-radio>
              <el-radio label="follow">关注</el-radio>
              <el-radio label="friend">好友</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="执行次数">
            <el-input-number v-model="formData.count" :min="1" :max="100" style="width:150px" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="executing" @click="handleExecute">执行</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <el-card v-if="resultData !== null" class="box-card mb-4">
        <template #header>
          <span>最近一次执行结果</span>
        </template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="记录ID">
            <el-tag>{{ resultData.recordId }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="成功次数">
            <el-tag type="success" size="large">{{ resultData.successCount }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="操作类型">
            <el-tag>{{ operationLabel }}</el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <el-card class="box-card">
        <template #header>
          <div class="flex justify-between items-center">
            <span>执行记录</span>
            <div>
              <el-select v-model="recordFilter.operation" placeholder="操作类型" clearable style="width:140px" class="mr-2">
                <el-option label="粉丝" value="fans" />
                <el-option label="关注" value="follow" />
                <el-option label="好友" value="friend" />
              </el-select>
              <el-select v-model="recordFilter.environmentKey" placeholder="按环境过滤" clearable filterable style="width:200px" class="mr-2">
                <el-option v-for="env in environmentOptions" :key="env.key" :label="env.name" :value="env.key" />
              </el-select>
              <el-button icon="refresh" @click="loadRecords">刷新</el-button>
            </div>
          </div>
        </template>
        <el-table :data="recordList" row-key="ID" style="width:100%">
          <el-table-column align="left" label="记录ID" prop="ID" width="80" />
          <el-table-column align="left" label="环境Key" prop="environmentKey" width="160" show-overflow-tooltip />
          <el-table-column align="left" label="目标用户ID" prop="userId" width="120" />
          <el-table-column align="left" label="操作" prop="operation" width="100">
            <template #default="scope">
              <el-tag>{{ operationLabelOf(scope.row.operation) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column align="left" label="数量" prop="count" width="80" />
          <el-table-column align="left" label="成功数" prop="successCount" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.successCount >= scope.row.count ? 'success' : 'warning'">
                {{ scope.row.successCount }} / {{ scope.row.count }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column align="left" label="状态" prop="status" width="100">
            <template #default="scope">
              <el-tag :type="recordStatusType(scope.row.status)">{{ recordStatusLabel(scope.row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column align="left" label="创建时间" prop="CreatedAt" width="180">
            <template #default="scope">
              <span>{{ formatDate(scope.row.CreatedAt) }}</span>
            </template>
          </el-table-column>
        </el-table>
        <div class="gva-pagination">
          <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="recordPage"
            :page-size="recordPageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="recordTotal"
            @current-change="(v) => { recordPage = v; loadRecords() }"
            @size-change="(v) => { recordPageSize = v; recordPage = 1; loadRecords() }"
          />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { createFanFollow, getFanFollowList } from '@/plugin/tools/api/fanFollow'
import { getEnvironmentList } from '@/plugin/tools/api/environment'
import { formatDate } from '@/utils/format'
import { ElMessage } from 'element-plus'
import { ref, computed, onMounted, onUnmounted } from 'vue'

defineOptions({ name: 'ToolsFanFollow' })

const formData = ref({
  environmentKey: '',
  userId: 0,
  operation: 'fans',
  count: 1
})

const environmentOptions = ref([])
const executing = ref(false)
const resultData = ref(null)

const recordList = ref([])
const recordTotal = ref(0)
const recordPage = ref(1)
const recordPageSize = ref(10)
const recordFilter = ref({ environmentKey: '', operation: '' })

let pollTimer = null

const operationLabel = computed(() => {
  const map = { fans: '粉丝', follow: '关注', friend: '好友' }
  return map[formData.value.operation] || formData.value.operation
})

const operationLabelOf = (op) => {
  const map = { fans: '粉丝', follow: '关注', friend: '好友' }
  return map[op] || op
}

const recordStatusType = (s) => {
  const map = { running: 'warning', completed: 'success', stopped: 'info', failed: 'danger' }
  return map[s] || ''
}

const recordStatusLabel = (s) => {
  const map = { running: '运行中', completed: '已完成', stopped: '已停止', failed: '失败' }
  return map[s] || s
}

const onEnvChange = () => {}

const loadEnvironments = async () => {
  const res = await getEnvironmentList({ page: 1, pageSize: 1000 })
  if (res.code === 0) {
    environmentOptions.value = res.data.list
  }
}

const loadRecords = async () => {
  const res = await getFanFollowList({
    page: recordPage.value,
    pageSize: recordPageSize.value,
    environmentKey: recordFilter.value.environmentKey || undefined,
    operation: recordFilter.value.operation || undefined
  })
  if (res.code === 0) {
    recordList.value = res.data.list
    recordTotal.value = res.data.total
  }
}

const startPolling = () => {
  stopPolling()
  pollTimer = setInterval(() => {
    if (recordList.value.some(r => r.status === 'running')) {
      loadRecords()
    }
  }, 5000)
}

const stopPolling = () => {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

onMounted(() => {
  loadEnvironments()
  loadRecords()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})

const handleExecute = async () => {
  if (!formData.value.environmentKey) {
    ElMessage.warning('请选择环境')
    return
  }
  if (!formData.value.userId) {
    ElMessage.warning('请输入用户ID')
    return
  }
  if (formData.value.count <= 0) {
    ElMessage.warning('执行次数必须大于0')
    return
  }
  executing.value = true
  resultData.value = null
  try {
    const res = await createFanFollow(formData.value)
    if (res.code === 0) {
      resultData.value = res.data
      ElMessage.success('执行完成')
      recordPage.value = 1
      loadRecords()
    } else {
      ElMessage.error(res.msg || '执行失败')
    }
  } finally {
    executing.value = false
  }
}
</script>