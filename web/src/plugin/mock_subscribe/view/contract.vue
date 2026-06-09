<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        :rules="searchRule"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="创建日期" prop="createdAt">
          <template #label>
            <span>
              创建日期
              <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker
            v-model="searchInfo.startCreatedAt"
            type="datetime"
            placeholder="开始日期"
            :disabled-date="(time) => searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"
          />
          —
          <el-date-picker
            v-model="searchInfo.endCreatedAt"
            type="datetime"
            placeholder="结束日期"
            :disabled-date="(time) => searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"
          />
        </el-form-item>
        <el-form-item label="用户OpenID" prop="openId">
          <el-input v-model="searchInfo.openId" placeholder="用户OpenID" clearable />
        </el-form-item>
        <el-form-item label="协议状态" prop="contractStatus">
          <el-select v-model="searchInfo.contractStatus" placeholder="协议状态" clearable style="width: 100%">
            <el-option label="签约成功" value="ACTIVE" />
            <el-option label="签约中" value="PENDING" />
            <el-option label="签约失败" value="FAILED" />
            <el-option label="已解约" value="TERMINATED" />
            <el-option label="已到期" value="EXPIRED" />
            <el-option label="已暂停" value="PAUSE" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        :row-key="getRowKey"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="日期" width="180">
          <template #default="scope">
            <span>{{ formatDateTime(scope.row.contract?.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="协议号" width="180" show-overflow-tooltip>
          <template #default="scope">
            <span>{{ scope.row.contract?.contractId || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="外部签约单号" width="180" show-overflow-tooltip>
          <template #default="scope">
            <span>{{ scope.row.contract?.outContractId || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="OpenID" width="180" show-overflow-tooltip>
          <template #default="scope">
            <span>{{ scope.row.contract?.openId || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="外部用户ID" width="180" show-overflow-tooltip>
          <template #default="scope">
            <span>{{ scope.row.contract?.outUserId || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="协议状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status?.contractStatus)">{{ getStatusLabel(scope.row.status?.contractStatus) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="签约流水号" width="220" show-overflow-tooltip>
          <template #default="scope">
            <span>{{ scope.row.status?.signSerialNo || scope.row.contract?.signSerialNo || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="解约方式" width="140" show-overflow-tooltip>
          <template #default="scope">
            <span>{{ getTerminateTypeLabel(scope.row.status?.terminateType) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="到期时间" width="180">
          <template #default="scope">
            <span>{{ formatDateTime(scope.row.status?.expireTime || scope.row.contract?.expireTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="首次扣款" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status?.isFirstDeduct ? 'success' : 'info'">{{ scope.row.status?.isFirstDeduct ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="已预通知" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status?.preNotifyCalled ? 'warning' : 'info'">{{ scope.row.status?.preNotifyCalled ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="最近预通知时间" width="180">
          <template #default="scope">
            <span>{{ formatDateTime(scope.row.status?.lastPreNotifyTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="回调通知地址" min-width="220" show-overflow-tooltip>
          <template #default="scope">
            <span>{{ scope.row.contract?.notifyUrl || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="160">
          <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button" @click="updateInfoFunc(scope.row)">变更状态</el-button>
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
          <span class="text-lg">变更协议状态</span>
          <div>
            <el-button type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>
      <el-form ref="elFormRef" :model="formData" label-position="top" :rules="rule" label-width="120px">
        <el-form-item label="协议号">
          <el-input :model-value="formData.contractId" disabled />
        </el-form-item>
        <el-form-item label="外部签约单号">
          <el-input :model-value="formData.outContractId" disabled />
        </el-form-item>
        <el-form-item label="协议状态" prop="contractStatus">
          <el-select v-model="formData.contractStatus" placeholder="选择协议状态" style="width: 100%">
            <el-option label="签约成功" value="ACTIVE" />
            <el-option label="签约中" value="PENDING" />
            <el-option label="签约失败" value="FAILED" />
            <el-option label="已解约" value="TERMINATED" />
            <el-option label="已到期" value="EXPIRED" />
            <el-option label="已暂停" value="PAUSE" />
          </el-select>
        </el-form-item>
        <el-form-item label="解约方式" prop="terminateType">
          <el-select v-model="formData.terminateType" placeholder="选择解约方式" style="width: 100%" :disabled="formData.contractStatus !== 'TERMINATED'">
            <el-option label="用户申请解约" value="USER_REQUEST" />
            <el-option label="商户申请解约" value="MERCHANT_REQUEST" />
            <el-option label="到期自动解约" value="EXPIRED" />
            <el-option label="系统撤销" value="SYSTEM_REVOKE" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { getContractList, findContract, updateContractStatus } from '@/plugin/mock_subscribe/api/contract'
import { formatDate } from '@/utils/format'
import { ElMessage } from 'element-plus'
import { reactive, ref, watch } from 'vue'

defineOptions({ name: 'MockSubscribeContract' })

const formData = ref({
  id: 0,
  contractId: '',
  outContractId: '',
  contractStatus: '',
  terminateType: ''
})

const rule = {}

const searchInfo = ref({
  startCreatedAt: null,
  endCreatedAt: null,
  openId: '',
  contractStatus: ''
})

const searchRule = reactive({})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

const elSearchFormRef = ref(null)
const elFormRef = ref(null)
const dialogFormVisible = ref(false)

const formatDateTime = (value) => {
  if (!value) return '-'
  const dateValue = value instanceof Date ? value : new Date(value)
  if (Number.isNaN(dateValue.getTime())) return '-'
  return formatDate(dateValue)
}

const getRowKey = (row) => row.contract?.ID ?? row.contract?.id ?? row.ID ?? row.id

const getStatusType = (status) => {
  switch (status) {
    case 'ACTIVE': return 'success'
    case 'PENDING': return 'warning'
    case 'FAILED': return 'danger'
    case 'TERMINATED': return 'info'
    case 'EXPIRED': return 'info'
    case 'PAUSE': return 'warning'
    default: return 'info'
  }
}

const getStatusLabel = (status) => {
  switch (status) {
    case 'ACTIVE': return '签约成功'
    case 'PENDING': return '签约中'
    case 'FAILED': return '签约失败'
    case 'TERMINATED': return '已解约'
    case 'EXPIRED': return '已到期'
    case 'PAUSE': return '已暂停'
    default: return status || '-'
  }
}

const getTerminateTypeLabel = (terminateType) => {
  switch (terminateType) {
    case 'USER_REQUEST': return '用户申请解约'
    case 'MERCHANT_REQUEST': return '商户申请解约'
    case 'EXPIRED': return '到期自动解约'
    case 'SYSTEM_REVOKE': return '系统撤销'
    default: return terminateType || '-'
  }
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
  const res = await getContractList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total || 0
  }
}

getTableData()

const onReset = () => {
  searchInfo.value = { startCreatedAt: null, endCreatedAt: null, openId: '', contractStatus: '' }
  getTableData()
}

const onSubmit = () => {
  elSearchFormRef.value?.validate((valid) => {
    if (valid) {
      page.value = 1
      getTableData()
    }
  })
}

watch(() => formData.value.contractStatus, (newVal) => {
  if (newVal !== 'TERMINATED') {
    formData.value.terminateType = ''
  }
})

const updateInfoFunc = async (row) => {
  const contractId = row.contract?.ID ?? row.contract?.id ?? row.ID ?? row.id
  const res = await findContract({ ID: contractId })
  if (res.code === 0) {
    formData.value = {
      id: res.data.contract?.ID ?? res.data.contract?.id ?? contractId,
      contractId: res.data.contract?.contractId ?? '',
      outContractId: res.data.contract?.outContractId ?? '',
      contractStatus: res.data.status?.contractStatus ?? '',
      terminateType: res.data.status?.terminateType ?? ''
    }
    dialogFormVisible.value = true
  }
}

const closeDialog = () => {
  dialogFormVisible.value = false
  elFormRef.value?.resetFields()
  formData.value = { id: 0, contractId: '', outContractId: '', contractStatus: '', terminateType: '' }
}

const enterDialog = async () => {
  const res = await updateContractStatus({
    id: formData.value.id,
    contractStatus: formData.value.contractStatus,
    terminateType: formData.value.terminateType
  })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '更新成功' })
    closeDialog()
    getTableData()
  }
}
</script>
