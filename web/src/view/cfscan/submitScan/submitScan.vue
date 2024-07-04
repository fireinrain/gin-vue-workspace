<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
      <el-form-item label="创建日期" prop="createdAt">
      <template #label>
        <span>
          创建日期
          <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
            <el-icon><QuestionFilled /></el-icon>
          </el-tooltip>
        </span>
      </template>
      <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始日期" :disabled-date="time=> searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"></el-date-picker>
       —
      <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束日期" :disabled-date="time=> searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"></el-date-picker>
      </el-form-item>
      
        <el-form-item label="扫描简介" prop="scanDesc">
         <el-input v-model="searchInfo.scanDesc" placeholder="搜索条件" />

        </el-form-item>
           <el-form-item label="扫描类型" prop="scanType">
            <el-select v-model="searchInfo.scanType" clearable placeholder="请选择" @clear="()=>{searchInfo.scanType=undefined}">
              <el-option v-for="(item,key) in CFScanTypeOptions" :key="key" :label="item.label" :value="item.value" />
            </el-select>
            </el-form-item>
        <el-form-item label="ASN编号" prop="asnNumber">
         <el-input v-model="searchInfo.asnNumber" placeholder="搜索条件" />

        </el-form-item>
           <el-form-item label="任务状态" prop="scanStatus">
            <el-select v-model="searchInfo.scanStatus" clearable placeholder="请选择" @clear="()=>{searchInfo.scanStatus=undefined}">
              <el-option v-for="(item,key) in ScanTaskStatusOptions" :key="key" :label="item.label" :value="item.value" />
            </el-select>
            </el-form-item>

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery=true" v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery=false" v-else>收起</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
            <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
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
        
        <el-table-column align="left" label="日期" prop="createdAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        
        <el-table-column align="left" label="扫描简介" prop="scanDesc" width="120" />
        <el-table-column align="left" label="扫描类型" prop="scanType" width="120">
            <template #default="scope">
            {{ filterDict(scope.row.scanType,CFScanTypeOptions) }}
            </template>
        </el-table-column>
        <el-table-column align="left" label="ASN编号" prop="asnNumber" width="120" />
        <el-table-column align="left" label="IP信息类型" prop="ipinfoType" width="120">
            <template #default="scope">
            {{ filterDict(scope.row.ipinfoType,IPInfoTypeOptions) }}
            </template>
        </el-table-column>
        <!--<el-table-column align="left" label="IP信息列表" prop="ipinfoList" width="120" />-->
        <el-table-column align="left" label="开启TLS" prop="enableTls" width="120">
            <template #default="scope">
            {{ filterDict(scope.row.enableTls,TLSScanTypeOptions) }}
            </template>
        </el-table-column>
        <el-table-column align="left" label="扫描端口" prop="scanPorts" width="120" />
        <el-table-column align="left" label="扫描速率" prop="scanRate" width="120" />
        <el-table-column align="left" label="IP检测线程" prop="ipcheckThread" width="120" />
        <el-table-column align="left" label="开启测速" prop="enableSpeedtest" width="120">
            <template #default="scope">
            {{ filterDict(scope.row.enableSpeedtest,IPSpeedTestOptions) }}
            </template>
        </el-table-column>
        <el-table-column align="left" label="任务状态" prop="scanStatus" width="120">
            <template #default="scope">
            {{ filterDict(scope.row.scanStatus,ScanTaskStatusOptions) }}
            </template>
        </el-table-column>
        <!--<el-table-column align="left" label="扫描结果" prop="scanResult" width="120" />-->
        <el-table-column align="left" label="操作" fixed="right" min-width="240">
            <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button" @click="updateSubmitScanFunc(scope.row)">查看</el-button>
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
    <el-drawer destroy-on-close size="800" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
       <template #header>
              <div class="flex justify-between items-center">
                <span class="text-lg">{{type==='create'?'添加':'修改'}}</span>
                <div>
                  <el-button type="primary" @click="enterDialog">确 定</el-button>
                  <el-button @click="closeDialog">取 消</el-button>
                </div>
              </div>
            </template>

          <el-form class="row-form" :model="formData" :disabled="formIsFormView" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
            <el-form-item label="扫描简介:"  prop="scanDesc" >
              <el-input v-model="formData.scanDesc" :clearable="true"  placeholder="请输入扫描简介" />
            </el-form-item>
            <el-form-item label="扫描类型:"  prop="scanType" >
              <el-select v-model="formData.scanType" placeholder="请选择扫描类型" style="width:100%" :clearable="true" @change="handleScanTypeChange">
                <el-option v-for="(item,key) in CFScanTypeOptions" :key="key" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="ASN编号:"  prop="asnNumber" >
              <el-input v-model="formData.asnNumber" :disabled="isAsnDisabled" :clearable="true"  placeholder="请输入ASN编号,多个ASN用,分割" />
            </el-form-item>
            <el-form-item label="IP类型:"  prop="ipinfoType" >
              <el-select v-model="formData.ipinfoType" :disabled="isIpDisabled" placeholder="请选择IP类型" style="width:100%" :clearable="true" >
                <el-option v-for="(item,key) in IPInfoTypeOptions" :key="key" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="IP信息列表:"  prop="ipinfoList" >
              <el-input :row="4" type="textarea" :disabled="isIpDisabled" v-model="formData.ipinfoList" :clearable="true"  placeholder="请输入IP信息列表" />
            </el-form-item>
            <el-form-item label="开启TLS:"  prop="enableTls" >
              <el-select v-model="formData.enableTls" placeholder="请选择开启TLS" style="width:100%" :clearable="true" >
                <el-option v-for="(item,key) in TLSScanTypeOptions" :key="key" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="扫描端口:"  prop="scanPorts" >
              <el-input v-model="formData.scanPorts" :clearable="true"  placeholder="请输入扫描端口" />
            </el-form-item>
            <el-form-item label="扫描速率:"  prop="scanRate" >
              <el-input v-model.number="formData.scanRate" :clearable="true" placeholder="请输入扫描速率" />
            </el-form-item>
            <el-form-item label="IP检查线程:"  prop="ipcheckThread" >
              <el-input v-model.number="formData.ipcheckThread" :clearable="true" placeholder="请输入IP检查线程" />
            </el-form-item>
            <el-form-item label="开启IP测速:"  prop="enableSpeedtest" >
              <el-select v-model="formData.enableSpeedtest" placeholder="请选择开启IP测速" style="width:100%" :clearable="true" >
                <el-option v-for="(item,key) in IPSpeedTestOptions" :key="key" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="任务状态:"  prop="scanStatus" >
              <el-select v-model="formData.scanStatus" placeholder="请选择任务状态" style="width:100%" :clearable="true" >
                <el-option v-for="(item,key) in ScanTaskStatusOptions" :key="key" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="任务结果:"  prop="scanResult" >
              <el-input :rows="16" type="textarea" v-model="formData.scanResult" :clearable="true"  placeholder="请输入任务结果" />
            </el-form-item>
          </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createSubmitScan,
  deleteSubmitScan,
  deleteSubmitScanByIds,
  findSubmitScan,
  getSubmitScanList,
  updateSubmitScan
} from '@/api/cfscan/submitScan'

// 全量引入格式化工具 请按需保留
import {filterDict, formatDate, getDictFunc} from '@/utils/format'
import {ElMessage, ElMessageBox} from 'element-plus'
import {reactive, ref} from 'vue'

defineOptions({
    name: 'SubmitScan'
})

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const TLSScanTypeOptions = ref([])
const IPSpeedTestOptions = ref([])
const ScanTaskStatusOptions = ref([])
const CFScanTypeOptions = ref([])
const IPInfoTypeOptions = ref([])
const formData = ref({
        scanDesc: '',
        scanType: '',
        asnNumber: '',
        ipinfoType: '',
        ipinfoList: '',
        enableTls: '1',
        scanPorts: '443',
        scanRate: 20000,
        ipcheckThread: 100,
        enableSpeedtest: '1',
        scanStatus: '',
        scanResult: '',
        })
// 控制记录详情查看
const formIsFormView = ref(false)


// 验证规则
const rule = reactive({
               scanDesc : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
               scanType : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
               enableTls : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
               scanPorts : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
               scanRate : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
              ],
               ipcheckThread : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
              ],
               enableSpeedtest : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
})
const isAsnDisabled = ref(false)
const isIpDisabled = ref(false)

const handleScanTypeChange = (value) => {
  isIpDisabled.value = false
  isAsnDisabled.value = false
  if(value === "1" || value === "2"){
    isIpDisabled.value = true
  }
  if(value === "3" || value === "4"){
    isAsnDisabled.value = true
  }
}



const searchRule = reactive({
  createdAt: [
    { validator: (rule, value, callback) => {
      if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change' }
  ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await getSubmitScanList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
    TLSScanTypeOptions.value = await getDictFunc('TLSScanType')
    IPSpeedTestOptions.value = await getDictFunc('IPSpeedTest')
    ScanTaskStatusOptions.value = await getDictFunc('ScanTaskStatus')
    CFScanTypeOptions.value = await getDictFunc('CFScanType')
    IPInfoTypeOptions.value = await getDictFunc('IPInfoType')
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteSubmitScanFunc(row)
        })
    }

// 多选删除
const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
      const IDs = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map(item => {
          IDs.push(item.ID)
        })
      const res = await deleteSubmitScanByIds({ IDs })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === IDs.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
      })
    }

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateSubmitScanFunc = async(row) => {
    const res = await findSubmitScan({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
    formIsFormView.value = true
}


// 删除行
const deleteSubmitScanFunc = async (row) => {
    const res = await deleteSubmitScan({ ID: row.ID })
    if (res.code === 0) {
        ElMessage({
                type: 'success',
                message: '删除成功'
            })
            if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        scanDesc: '',
        scanType: '',
        asnNumber: '',
        ipinfoType: '',
        ipinfoList: '',
        enableTls: '',
        scanPorts: '',
        scanRate: undefined,
        ipcheckThread: undefined,
        enableSpeedtest: '',
        scanStatus: '',
        scanResult: '',
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createSubmitScan(formData.value)
                  break
                case 'update':
                  res = await updateSubmitScan(formData.value)
                  break
                default:
                  res = await createSubmitScan(formData.value)
                  break
              }
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: '创建/更改成功'
                })
                closeDialog()
                getTableData()
              }
      })
}

</script>

<style>
.row-form {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
}

.row-form .el-form-item {
  width: 48%;
}

</style>
