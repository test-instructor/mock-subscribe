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
        <el-form-item label="外部签约单号" prop="outContractCode">
          <el-input v-model="searchInfo.outContractCode" placeholder="外部签约单号" clearable />
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
        <el-table-column align="left" label="回调类型" prop="callbackType" width="180" />
        <el-table-column align="left" label="商户号" prop="mchId" width="140" />
        <el-table-column align="left" label="协议号" prop="contractCode" width="180" show-overflow-tooltip />
        <el-table-column align="left" label="外部签约单号" prop="outContractCode" width="180" show-overflow-tooltip />
        <el-table-column align="left" label="签约号" prop="signSerialNo" width="180" show-overflow-tooltip />
        <el-table-column align="left" label="协议状态" prop="contractStatus" width="120" />
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
    <el-drawer v-model="drawerVisible" destroy-on-close size="760px" :show-close="false" title="回调记录详情">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="记录ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item label="回调类型">{{ detailData.callbackType }}</el-descriptions-item>
        <el-descriptions-item label="商户号">{{ detailData.mchId }}</el-descriptions-item>
        <el-descriptions-item label="来源IP">{{ detailData.sourceIp }}</el-descriptions-item>
        <el-descriptions-item label="协议号">{{ detailData.contractCode }}</el-descriptions-item>
        <el-descriptions-item label="外部签约单号">{{ detailData.outContractCode }}</el-descriptions-item>
        <el-descriptions-item label="签约号">{{ detailData.signSerialNo }}</el-descriptions-item>
        <el-descriptions-item label="协议状态">{{ detailData.contractStatus }}</el-descriptions-item>
        <el-descriptions-item label="签名">{{ detailData.sign }}</el-descriptions-item>
        <el-descriptions-item label="验签结果">{{ detailData.signValid ? '通过' : '失败' }}</el-descriptions-item>
        <el-descriptions-item label="错误信息" :span="2">{{ detailData.signErrorMessage }}</el-descriptions-item>
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
import { getCallbackRecordList, findCallbackRecord } from '@/plugin/mock_subscribe/api/callback'
import { formatDate } from '@/utils/format'
import { reactive, ref } from 'vue'

defineOptions({ name: 'MockSubscribeCallback' })

const searchInfo = ref({
  startCreatedAt: null,
  endCreatedAt: null,
  mchId: '',
  outContractCode: '',
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
  const res = await getCallbackRecordList({
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
      getTableData()
    }
  })
}

const onReset = () => {
  searchInfo.value = {
    startCreatedAt: null,
    endCreatedAt: null,
    mchId: '',
    outContractCode: '',
    signValid: undefined
  }
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

const viewDetail = async (row) => {
  const res = await findCallbackRecord({ ID: row.ID })
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
