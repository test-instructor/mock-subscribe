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
        <el-form-item label="名称" prop="name">
          <el-input v-model="searchInfo.name" placeholder="中文名称" clearable />
        </el-form-item>
        <el-form-item label="Key" prop="key">
          <el-input v-model="searchInfo.key" placeholder="唯一标识" clearable />
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
        <el-table-column align="left" label="名称" prop="name" width="150" show-overflow-tooltip />
        <el-table-column align="left" label="Key" prop="key" width="180" show-overflow-tooltip />
        <el-table-column align="left" label="服务器地址" prop="domain" width="200" show-overflow-tooltip />
        <el-table-column align="left" label="端口" prop="port" width="100" />
        <el-table-column align="left" label="备注" prop="remark" min-width="200" show-overflow-tooltip />
        <el-table-column align="left" label="操作" fixed="right" width="180">
          <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button" @click="updateInfoFunc(scope.row)">编辑</el-button>
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
          <span class="text-lg">{{ type === 'create' ? '添加' : '修改' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>
      <el-form ref="elFormRef" :model="formData" label-position="top" :rules="rule" label-width="80px">
        <el-form-item label="中文名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入中文名称" clearable />
        </el-form-item>
        <el-form-item label="唯一标识(Key)" prop="key">
          <el-input v-model="formData.key" placeholder="请输入唯一标识" :disabled="type === 'update'" clearable />
        </el-form-item>
        <el-form-item label="服务器地址" prop="domain">
          <el-input v-model="formData.domain" placeholder="请输入服务器地址(含http协议)" clearable />
        </el-form-item>
        <el-form-item label="端口号" prop="port">
          <el-input-number v-model="formData.port" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入备注" clearable />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { getEnvironmentList, createEnvironment, updateEnvironment, deleteEnvironment, findEnvironment } from '@/plugin/tools/api/environment'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({ name: 'ToolsEnvironment' })

const formData = ref({
  name: '',
  key: '',
  domain: '',
  port: 80,
  remark: ''
})

const rule = {
  name: [{ required: true, message: '请输入中文名称', trigger: 'blur' }],
  key: [{ required: true, message: '请输入唯一标识', trigger: 'blur' }],
  domain: [{ required: true, message: '请输入服务器地址', trigger: 'blur' }],
}

const searchInfo = ref({
  name: '',
  key: ''
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

const onReset = () => {
  searchInfo.value = { name: '', key: '' }
  getTableData()
}

const onSubmit = () => {
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
  const res = await getEnvironmentList({
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
  const res = await findEnvironment({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data
    type.value = 'update'
    dialogFormVisible.value = true
  }
}

const closeDialog = () => {
  dialogFormVisible.value = false
  elFormRef.value?.resetFields()
  formData.value = { name: '', key: '', domain: '', port: 80, remark: '' }
}

const deleteRow = async (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    .then(async () => {
      const res = await deleteEnvironment({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '删除成功' })
        getTableData()
      }
    })
}

const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (valid) {
      let res
      if (type.value === 'create') {
        res = await createEnvironment(formData.value)
      } else {
        res = await updateEnvironment(formData.value)
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
