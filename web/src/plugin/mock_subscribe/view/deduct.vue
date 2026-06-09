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
        <el-form-item label="操作类型" prop="operationType">
          <el-select v-model="searchInfo.operationType" placeholder="操作类型" clearable style="width: 100%">
            <el-option label="签约" value="sign" />
            <el-option label="查询签约" value="query" />
            <el-option label="解约" value="terminate" />
            <el-option label="扣款" value="deduct" />
            <el-option label="预扣费通知" value="pre_notify" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="searchInfo.status" placeholder="状态" clearable style="width: 100%">
            <el-option label="成功" value="SUCCESS" />
            <el-option label="失败" value="FAILED" />
            <el-option label="进行中" value="PENDING" />
            <el-option label="退款中" value="REFUNDING" />
            <el-option label="已退款" value="REFUNDED" />
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
        row-key="ID"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="日期" prop="CreatedAt" width="180">
          <template #default="scope">
            <span>{{ formatDate(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作类型" prop="operationType" width="140">
          <template #default="scope">
            <el-tag :type="getOpType(scope.row.operationType)">{{ getOpLabel(scope.row.operationType) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="商户订单号" prop="outTradeNo" width="180" show-overflow-tooltip />
        <el-table-column align="left" label="微信交易单号" prop="transactionId" width="200" show-overflow-tooltip />
        <el-table-column align="left" label="金额(分)" prop="amount" width="100">
          <template #default="scope">
            <span>{{ scope.row.amount ? scope.row.amount / 100 + '元' : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" prop="status" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">{{ scope.row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="首次扣款" prop="isFirstDeduct" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.isFirstDeduct ? 'success' : 'info'">{{ scope.row.isFirstDeduct ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="已预通知" prop="preNotifyCalled" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.preNotifyCalled ? 'warning' : 'info'">{{ scope.row.preNotifyCalled ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="错误码" prop="errorCode" width="140" show-overflow-tooltip />
        <el-table-column align="left" label="错误信息" prop="errorMessage" width="200" show-overflow-tooltip />
        <el-table-column align="left" label="操作" fixed="right" min-width="120">
          <template #default="scope">
            <el-button type="primary" link icon="view" class="table-button" @click="viewDetail(scope.row)">详情</el-button>
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
      size="700px"
      :show-close="false"
      title="扣款记录详情"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="记录ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item label="操作类型">{{ getOpLabel(detailData.operationType) }}</el-descriptions-item>
        <el-descriptions-item label="商户订单号">{{ detailData.outTradeNo }}</el-descriptions-item>
        <el-descriptions-item label="微信交易单号">{{ detailData.transactionId }}</el-descriptions-item>
        <el-descriptions-item label="金额(分)">{{ detailData.amount }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ detailData.status }}</el-descriptions-item>
        <el-descriptions-item label="首次扣款">{{ detailData.isFirstDeduct ? '是' : '否' }}</el-descriptions-item>
        <el-descriptions-item label="已预通知">{{ detailData.preNotifyCalled ? '是' : '否' }}</el-descriptions-item>
        <el-descriptions-item label="错误码">{{ detailData.errorCode }}</el-descriptions-item>
        <el-descriptions-item label="错误信息">{{ detailData.errorMessage }}</el-descriptions-item>
        <el-descriptions-item label="回调地址" :span="2">{{ detailData.callbackUrl }}</el-descriptions-item>
        <el-descriptions-item label="回调结果" :span="2">{{ detailData.callbackResult }}</el-descriptions-item>
        <el-descriptions-item label="回调时间" :span="2">{{ detailData.callbackTime ? formatDate(new Date(detailData.callbackTime * 1000)) : '-' }}</el-descriptions-item>
        <el-descriptions-item label="请求数据" :span="2">
          <pre style="max-height: 200px; overflow: auto; white-space: pre-wrap; word-break: break-all;">{{ detailData.requestData }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="响应数据" :span="2">
          <pre style="max-height: 200px; overflow: auto; white-space: pre-wrap; word-break: break-all;">{{ detailData.responseData }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import { getDeductRecordList, findDeductRecord } from '@/plugin/mock_subscribe/api/deduct'
import { formatDate } from '@/utils/format'
import { ref, reactive } from 'vue'

defineOptions({ name: 'MockSubscribeDeduct' })

const searchInfo = ref({
  startCreatedAt: null,
  endCreatedAt: null,
  operationType: '',
  status: ''
})

const searchRule = reactive({})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

const elSearchFormRef = ref(null)

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const getTableData = async () => {
  const res = await getDeductRecordList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}

getTableData()

const onReset = () => {
  searchInfo.value = { startCreatedAt: null, endCreatedAt: null, operationType: '', status: '' }
  getTableData()
}

const onSubmit = () => {
  elSearchFormRef.value?.validate((valid) => {
    if (valid) {
      getTableData()
    }
  })
}

const getOpType = (op) => {
  switch (op) {
    case 'sign': return 'success'
    case 'query': return ''
    case 'terminate': return 'warning'
    case 'deduct': return 'primary'
    case 'pre_notify': return 'info'
    default: return 'info'
  }
}

const getOpLabel = (op) => {
  switch (op) {
    case 'sign': return '签约'
    case 'query': return '查询签约'
    case 'terminate': return '解约'
    case 'deduct': return '扣款'
    case 'pre_notify': return '预扣费通知'
    default: return op
  }
}

const getStatusType = (status) => {
  switch (status) {
    case 'SUCCESS': return 'success'
    case 'FAILED': return 'danger'
    case 'PENDING': return 'warning'
    case 'REFUNDING': return 'warning'
    case 'REFUNDED': return 'info'
    default: return 'info'
  }
}

const dialogFormVisible = ref(false)
const detailData = ref({})

const viewDetail = async (row) => {
  const res = await findDeductRecord({ ID: row.ID })
  if (res.code === 0) {
    detailData.value = res.data
    dialogFormVisible.value = true
  }
}
</script>
