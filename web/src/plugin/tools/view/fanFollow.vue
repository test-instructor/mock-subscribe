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

      <el-card v-if="resultData !== null" class="box-card">
        <template #header>
          <span>执行结果</span>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="成功次数">
            <el-tag type="success" size="large">{{ resultData.successCount }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="操作类型">
            <el-tag>{{ operationLabel }}</el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { createFanFollow } from '@/plugin/tools/api/fanFollow'
import { getEnvironmentList } from '@/plugin/tools/api/environment'
import { ElMessage } from 'element-plus'
import { ref, computed, onMounted } from 'vue'

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

const operationLabel = computed(() => {
  const map = { fans: '粉丝', follow: '关注', friend: '好友' }
  return map[formData.value.operation] || formData.value.operation
})

const onEnvChange = () => {
}

const loadEnvironments = async () => {
  const res = await getEnvironmentList({ page: 1, pageSize: 1000 })
  if (res.code === 0) {
    environmentOptions.value = res.data.list
  }
}

onMounted(() => {
  loadEnvironments()
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
    } else {
      ElMessage.error(res.msg || '执行失败')
    }
  } finally {
    executing.value = false
  }
}
</script>
