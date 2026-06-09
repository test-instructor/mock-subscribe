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
        <el-form-item label="应用ID" prop="appId">
          <el-input v-model="searchInfo.appId" placeholder="应用ID" clearable />
        </el-form-item>
        <el-form-item label="商户号" prop="mchId">
          <el-input v-model="searchInfo.mchId" placeholder="商户号" clearable />
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
        <el-button icon="delete" style="margin-left: 10px" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="日期" prop="CreatedAt" width="180">
          <template #default="scope">
            <span>{{ formatDate(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="应用ID" prop="appId" width="140" show-overflow-tooltip />
        <el-table-column align="left" label="商户号" prop="mchId" width="120" show-overflow-tooltip />
        <el-table-column align="left" label="签约AppID" prop="contractAppId" width="140" show-overflow-tooltip />
        <el-table-column align="left" label="签约模板" prop="contractTemplate" width="120" show-overflow-tooltip />
        <el-table-column align="left" label="验签" prop="verifySign" width="90">
          <template #default="scope">
            <el-tag :type="scope.row.verifySign ? 'success' : 'info'">
              {{ scope.row.verifySign ? '开启' : '关闭' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="签约回调" prop="signCallbackEnabled" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.signCallbackEnabled ? 'success' : 'info'">
              {{ scope.row.signCallbackEnabled ? '开启' : '关闭' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="扣款回调" prop="deductCallbackEnabled" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.deductCallbackEnabled ? 'success' : 'info'">
              {{ scope.row.deductCallbackEnabled ? '开启' : '关闭' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="签约目标状态" prop="signTargetStatus" width="120" />
        <el-table-column align="left" label="签约延时(秒)" prop="signStatusDelay" width="110" />
        <el-table-column align="left" label="扣款目标状态" prop="deductTargetStatus" width="120" />
        <el-table-column align="left" label="扣款延时(秒)" prop="deductStatusDelay" width="110" />
        <el-table-column align="left" label="签约时长(分)" prop="signDurationMinutes" width="110" />
        <el-table-column align="left" label="启用" prop="active" width="80">
          <template #default="scope">
            <el-switch v-model="scope.row.active" disabled />
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="200">
          <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button" @click="updateInfoFunc(scope.row)">变更</el-button>
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
      size="900px"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '添加' : '修改' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>
      <el-form ref="elFormRef" :model="formData" label-position="top" :rules="rule" label-width="80px">
        <el-form-item label="应用ID" prop="appId">
          <el-input v-model="formData.appId" placeholder="请输入应用ID" clearable />
        </el-form-item>
        <el-form-item label="商户号" prop="mchId">
          <el-input v-model="formData.mchId" placeholder="请输入商户号" clearable />
        </el-form-item>
        <el-form-item label="签约商户号" prop="contractMchId">
          <el-input v-model="formData.contractMchId" placeholder="请输入签约商户号" clearable />
        </el-form-item>
        <el-form-item label="签约AppID" prop="contractAppId">
          <el-input v-model="formData.contractAppId" placeholder="请输入签约AppID" clearable />
        </el-form-item>
        <el-form-item label="用户账户展示名称" prop="displayName">
          <el-input v-model="formData.displayName" placeholder="请输入用户账户展示名称" clearable />
        </el-form-item>
        <el-form-item label="签名Key" prop="signKey">
          <el-input v-model="formData.signKey" placeholder="请输入签名Key" clearable />
        </el-form-item>
        <el-form-item label="是否验签" prop="verifySign">
          <el-switch v-model="formData.verifySign" />
        </el-form-item>
        <el-form-item label="签约模板" prop="contractTemplate">
          <el-input v-model="formData.contractTemplate" placeholder="请输入签约模板" clearable />
        </el-form-item>
        <el-divider content-position="left">签约配置</el-divider>
        <div class="flex flex-wrap gap-4">
          <el-form-item label="签约回调开关" prop="signCallbackEnabled">
            <el-switch v-model="formData.signCallbackEnabled" />
          </el-form-item>
          <el-form-item label="签约回调延时(秒)" prop="signCallbackDelay">
            <el-input-number v-model="formData.signCallbackDelay" :min="0" :max="3600" />
          </el-form-item>
          <el-form-item label="签约目标状态" prop="signTargetStatus">
            <el-select v-model="formData.signTargetStatus" placeholder="选择状态" style="width: 100%" clearable>
              <el-option label="签约成功(ACTIVE)" value="ACTIVE" />
              <el-option label="签约失败(FAILED)" value="FAILED" />
              <el-option label="签约中(PENDING)" value="PENDING" />
            </el-select>
          </el-form-item>
          <el-form-item label="签约状态延时(秒)" prop="signStatusDelay">
            <el-input-number v-model="formData.signStatusDelay" :min="0" :max="3600" />
          </el-form-item>
          <el-form-item label="签约时长(分钟)" prop="signDurationMinutes">
            <el-input-number v-model="formData.signDurationMinutes" :min="1" :max="43200" />
          </el-form-item>
        </div>
        <el-divider content-position="left">扣款配置</el-divider>
        <div class="flex flex-wrap gap-4">
          <el-form-item label="扣款回调开关" prop="deductCallbackEnabled">
            <el-switch v-model="formData.deductCallbackEnabled" />
          </el-form-item>
          <el-form-item label="扣款回调延时(秒)" prop="deductCallbackDelay">
            <el-input-number v-model="formData.deductCallbackDelay" :min="0" :max="3600" />
          </el-form-item>
          <el-form-item label="扣款目标状态" prop="deductTargetStatus">
            <el-select v-model="formData.deductTargetStatus" placeholder="选择状态" style="width: 100%" clearable>
              <el-option label="扣款成功(SUCCESS)" value="SUCCESS" />
              <el-option label="扣款失败(FAILED)" value="FAILED" />
              <el-option label="扣款中(PENDING)" value="PENDING" />
              <el-option label="退款中(REFUNDING)" value="REFUNDING" />
              <el-option label="已退款(REFUNDED)" value="REFUNDED" />
            </el-select>
          </el-form-item>
          <el-form-item label="扣款状态延时(秒)" prop="deductStatusDelay">
            <el-input-number v-model="formData.deductStatusDelay" :min="0" :max="3600" />
          </el-form-item>
          <el-form-item label="解约通知商户" prop="terminateNotifyEnabled">
            <el-switch v-model="formData.terminateNotifyEnabled" />
          </el-form-item>
          <el-form-item label="严格扣费规则" prop="strictDeductRule">
            <el-switch v-model="formData.strictDeductRule" />
          </el-form-item>
        </div>
        <el-divider content-position="left">启用状态</el-divider>
        <el-form-item label="启用" prop="active">
          <el-switch v-model="formData.active" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { getMerchantList, createMerchant, updateMerchant, deleteMerchant, findMerchant } from '@/plugin/mock_subscribe/api/merchant'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({ name: 'MockSubscribeMerchant' })

const formData = ref({
  appId: '',
  mchId: '',
  contractMchId: '',
  contractAppId: '',
  displayName: '',
  signKey: '',
  verifySign: true,
  contractTemplate: '',
  signCallbackEnabled: true,
  signCallbackDelay: 0,
  deductCallbackEnabled: true,
  deductCallbackDelay: 0,
  signTargetStatus: 'ACTIVE',
  signStatusDelay: 0,
  deductTargetStatus: 'SUCCESS',
  deductStatusDelay: 0,
  terminateNotifyEnabled: false,
  signDurationMinutes: 1440,
  strictDeductRule: false,
  active: true
})

const rule = {
  appId: [{ required: true, message: '请输入应用ID', trigger: 'blur' }],
  mchId: [{ required: true, message: '请输入商户号', trigger: 'blur' }],
  signKey: [{ validator: (_, value, callback) => {
    if (formData.value.verifySign && !value) {
      callback(new Error('验签开启时请输入签名Key'))
      return
    }
    callback()
  }, trigger: 'blur' }]
}

const searchInfo = ref({
  startCreatedAt: null,
  endCreatedAt: null,
  appId: '',
  mchId: ''
})

const searchRule = reactive({})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

const multipleTable = ref(null)
const multipleSelection = ref([])

const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const onReset = () => {
  searchInfo.value = { startCreatedAt: null, endCreatedAt: null, appId: '', mchId: '' }
  getTableData()
}

const onSubmit = () => {
  elSearchFormRef.value?.validate((valid) => {
    if (valid) {
      getTableData()
    }
  })
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
  const res = await getMerchantList({
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

const dialogFormVisible = ref(false)
const type = ref('')

const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}

const updateInfoFunc = async (row) => {
  const res = await findMerchant({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data
    type.value = 'update'
    dialogFormVisible.value = true
  }
}

const closeDialog = () => {
  dialogFormVisible.value = false
  elFormRef.value?.resetFields()
  formData.value = {
    appId: '', mchId: '', contractMchId: '', contractAppId: '', displayName: '',
    signKey: '', verifySign: true, contractTemplate: '', signCallbackEnabled: true, signCallbackDelay: 0,
    deductCallbackEnabled: true, deductCallbackDelay: 0, signTargetStatus: 'ACTIVE',
    signStatusDelay: 0, deductTargetStatus: 'SUCCESS', deductStatusDelay: 0,
    terminateNotifyEnabled: false, signDurationMinutes: 1440, strictDeductRule: false, active: true
  }
}

const deleteRow = async (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    .then(async () => {
      const res = await deleteMerchant({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '删除成功' })
        getTableData()
      }
    })
}

const onDelete = async () => {
  ElMessageBox.confirm('确定要删除选中的记录吗?', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    .then(async () => {
      const IDs = multipleSelection.value.map((item) => item.ID)
      for (const id of IDs) {
        await deleteMerchant({ ID: id })
      }
      ElMessage({ type: 'success', message: '批量删除成功' })
      getTableData()
    })
}

const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (valid) {
      let res
      if (type.value === 'create') {
        res = await createMerchant(formData.value)
      } else {
        res = await updateMerchant(formData.value)
      }
      if (res.code === 0) {
        ElMessage({ type: 'success', message: type.value === 'create' ? '创建成功' : '更新成功' })
        closeDialog()
        getTableData()
      }
    }
  })
}
</script>
