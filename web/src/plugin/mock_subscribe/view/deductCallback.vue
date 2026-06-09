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
          <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始日期" />
          —
          <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束日期" />
        </el-form-item>
        <el-form-item label="商户号" prop="mchId">
          <el-input v-model="searchInfo.mchId" placeholder="商户号" clearable />
        </el-form-item>
        <el-form-item label="商户订单号" prop="outTradeNo">
          <el-input v-model="searchInfo.outTradeNo" placeholder="商户订单号" clearable />
        </el-form-item>
        <el-form-item label="微信交易单号" prop="transactionId">
          <el-input v-model="searchInfo.transactionId" placeholder="微信交易单号" clearable />
        </el-form-item>
        <el-form-item label="交易状态" prop="tradeState">
          <el-select v-model="searchInfo.tradeState" placeholder="交易状态" clearable style="width: 140px">
            <el-option label="成功" value="SUCCESS" />
            <el-option label="失败" value="FAILED" />
            <el-option label="进行中" value="PENDING" />
            <el-option label="退款中" value="REFUNDING" />
            <el-option label="已退款" value="REFUNDED" />
          </el-select>
        </el-form-item>
        <el-form-item label="验签结果" prop="signValid">
          <el-select v-model="searchInfo.signValid" placeholder="验签结果" clearable style="width: 120px">
            <el-option label="通过" :value="true" />
            <el-option label="失败" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table style="width: 100%" :data="tableData" row-key="ID" tooltip-effect="dark">
        <el-table-column align="left" label="时间" width="180">
          <template #default="scope">
            <span>{{ formatDate(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="商户号" prop="mchId" width="140" />
        <el-table-column align="left" label="商户订单号" prop="outTradeNo" width="180" show-overflow-tooltip />
        <el-table-column align="left" label="微信交易单号" prop="transactionId" width="200" show-overflow-tooltip />
        <el-table-column align="left" label="交易类型" prop="tradeType" width="100" />
        <el-table-column align="left" label="交易状态" prop="tradeState" width="120">
          <template #default="scope">
            <el-tag :type="getTradeStateType(scope.row.tradeState)">{{ scope.row.tradeState || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="金额(分)" prop="totalAmount" width="100" />
        <el-table-column align="left" label="来源IP" prop="sourceIp" width="140" />
        <el-table-column align="left" label="验签" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.signValid ? 'success' : 'danger'">
              {{ scope.row.signValid ? '通过' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="错误信息" prop="signErrorMessage" min-width="220" show-overflow-tooltip />
        <el-table-column align="left" label="操作" fixed="right" min-width="120">
          <template #default="scope">
            <el-button type="primary" link icon="view" @click="viewDetail(scope.row)">详情</el-button>
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
    <el-drawer v-model="drawerVisible" destroy-on-close size="760px" :show-close="false" title="代扣回调详情">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="记录ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item label="商户号">{{ detailData.mchId }}</el-descriptions-item>
        <el-descriptions-item label="商户订单号">{{ detailData.outTradeNo }}</el-descriptions-item>
        <el-descriptions-item label="微信交易单号">{{ detailData.transactionId }}</el-descriptions-item>
        <el-descriptions-item label="交易类型">{{ detailData.tradeType }}</el-descriptions-item>
        <el-descriptions-item label="交易状态">{{ detailData.tradeState }}</el-descriptions-item>
        <el-descriptions-item label="总金额(分)">{{ detailData.totalAmount }}</el-descriptions-item>
        <el-descriptions-item label="现金金额(分)">{{ detailData.cashAmount }}</el-descriptions-item>
        <el-descriptions-item label="支付完成时间">{{ detailData.timeEnd || '-' }}</el-descriptions-item>
        <el-descriptions-item label="来源IP">{{ detailData.sourceIp }}</el-descriptions-item>
        <el-descriptions-item label="关联签约ID">{{ detailData.contractIdRef }}</el-descriptions-item>
        <el-descriptions-item label="关联扣款记录ID">{{ detailData.deductRecordIdRef }}</el-descriptions-item>
        <el-descriptions-item label="签名">{{ detailData.sign }}</el-descriptions-item>
        <el-descriptions-item label="验签结果">{{ detailData.signValid ? '通过' : '失败' }}</el-descriptions-item>
        <el-descriptions-item label="验签错误" :span="2">{{ detailData.signErrorMessage }}</el-descriptions-item>
        <el-descriptions-item label="请求头" :span="2">
          <pre class="pre-block">{{ detailData.headers }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="原始请求体" :span="2">
          <pre class="pre-block">{{ detailData.rawBody }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="ACK" :span="2">
          <pre class="pre-block">{{ detailData.ackXml }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import { getDeductCallbackRecordList, findDeductCallbackRecord } from '@/plugin/mock_subscribe/api/deductCallback'
import { formatDate } from '@/utils/format'
import { reactive, ref } from 'vue'

defineOptions({ name: 'MockSubscribeDeductCallback' })

const searchInfo = ref({
  startCreatedAt: null,
  endCreatedAt: null,
  mchId: '',
  outTradeNo: '',
  transactionId: '',
  tradeState: '',
  signValid: undefined
})

const searchRule = reactive({})
const elSearchFormRef = ref(null)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const drawerVisible = ref(false)
const detailData = ref({})

const getTableData = async () => {
  const res = await getDeductCallbackRecordList({
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

const onSubmit = () => {
  elSearchFormRef.value?.validate((valid) => {
    if (valid) {
      page.value = 1
      getTableData()
    }
  })
}

const onReset = () => {
  searchInfo.value = {
    startCreatedAt: null,
    endCreatedAt: null,
    mchId: '',
    outTradeNo: '',
    transactionId: '',
    tradeState: '',
    signValid: undefined
  }
  page.value = 1
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const getTradeStateType = (status) => {
  switch (status) {
    case 'SUCCESS': return 'success'
    case 'FAILED': return 'danger'
    case 'PENDING': return 'warning'
    case 'REFUNDING': return 'warning'
    case 'REFUNDED': return 'info'
    default: return 'info'
  }
}

const viewDetail = async (row) => {
  const res = await findDeductCallbackRecord({ ID: row.ID })
  if (res.code === 0) {
    detailData.value = res.data
    drawerVisible.value = true
  }
}
</script>

<style scoped>
.pre-block {
  max-height: 220px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
