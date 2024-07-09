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
      
        <el-form-item label="任务描述" prop="taskDesc">
         <el-input v-model="searchInfo.taskDesc" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="ASN编号" prop="asnNumber">
         <el-input v-model="searchInfo.asnNumber" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="ASN描述" prop="asnDesc">
         <el-input v-model="searchInfo.asnDesc" placeholder="搜索条件" />

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
        @sort-change="sortChange"
        >
        <el-table-column type="selection" width="55" />
        

        <el-table-column sortable align="left" label="ID" prop="ID" width="90">
        </el-table-column>
        
        <el-table-column align="left" label="任务描述" prop="taskDesc" width="120" />
        <el-table-column sortable align="left" label="ASN编号" prop="asnNumber" width="120" />
        <el-table-column align="left" label="ASN描述" prop="asnDesc" width="120" />
        <el-table-column align="left" label="定时表达式" prop="crontabStr" width="120" />
        <!--<el-table-column align="left" label="任务配置" prop="taskConfig" width="120" />-->

        <el-table-column align="left" label="任务状态" prop="taskStatus" width="120">
            <template #default="scope">
            {{ filterDict(scope.row.taskStatus,ScheduleTaskStatusOptions) }}
            </template>
        </el-table-column>

        <!--<el-table-column align="left" label="是否开启" prop="enable" width="120">-->
        <!--  <template #default="scope">-->
        <!--    {{ filterDict(scope.row.enable,EnableOrNotOptions) }}-->
        <!--  </template>-->
        <!--</el-table-column>-->
          <el-table-column
              align="left"
              label="是否启用"
              min-width="120"
          >
            <template #default="scope">
              <el-switch
                  v-model="scope.row.enable"
                  inline-prompt
                  :active-value="'1'"
                  :inactive-value="'0'"
                  @change="()=>{switchEnable(scope.row)}"
              />
            </template>
          </el-table-column>
          <el-table-column align="left" label="创建日期" prop="createdAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
          </el-table-column>
          <el-table-column align="left" label="最后运行日期" prop="lastRunAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.lastRunAt) }}</template>
          </el-table-column>
          <el-table-column align="left" label="下次运行日期" prop="createdAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.nextRunAt) }}</template>
          </el-table-column>

        <el-table-column align="left" label="操作" fixed="right" min-width="240">
            <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button" @click="updateScheduleTaskFunc(scope.row)">变更</el-button>
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

          <el-form class="row-form" :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
            <el-form-item label="任务描述:"  prop="taskDesc" >
              <el-input v-model="formData.taskDesc" :clearable="true"  placeholder="请输入任务描述" />
            </el-form-item>
            <el-form-item label="ASN编号:"  prop="asnNumber" >
              <el-input v-model="formData.asnNumber" :clearable="true"  placeholder="请输入ASN编号" />
            </el-form-item>
            <el-form-item label="ASN描述:"  prop="asnDesc" >
              <el-input v-model="formData.asnDesc" :clearable="true"  placeholder="请输入ASN描述" />
            </el-form-item>
            <el-form-item label="定时表达式:"  prop="crontabStr" >
              <el-input v-model="formData.crontabStr" :clearable="true"  placeholder="分钟 小时 日期 月份 星期几,例如: 0 2 * * *" />
            </el-form-item>

            <el-form-item label="是否开启:"  prop="enable" >
              <el-select v-model="formData.enable" placeholder="请选择是否开启" style="width:100%" :clearable="true" >
                <el-option v-for="(item,key) in EnableOrNotOptions" :key="key" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="任务状态:"  prop="taskStatus" >
              <el-select v-model="formData.taskStatus" placeholder="请选择任务状态" style="width:100%" :clearable="true" >
                <el-option v-for="(item,key) in ScheduleTaskStatusOptions" :key="key" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="任务配置:"  prop="taskConfig" >
              <el-input type="textarea" :rows="12" v-model="formData.taskConfig" :clearable="true"  :placeholder=configSet />
            </el-form-item>
            <el-form-item label="配置模版(复制修改填入左边):" prop="taskConfigTemp" >
              <el-input type="textarea" :rows="12" disabled v-model="taskConfigTemp" :clearable="true" />
            </el-form-item>
          </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createScheduleTask,
  deleteScheduleTask,
  deleteScheduleTaskByIds,
  findScheduleTask,
  getScheduleTaskList,
  updateScheduleTask
} from '@/api/cfscan/scheduleTask'

// 全量引入格式化工具 请按需保留
import {filterDict, formatDate, getDictFunc} from '@/utils/format'
import {ElMessage, ElMessageBox} from 'element-plus'
import {nextTick, reactive, ref} from 'vue'

defineOptions({
    name: 'ScheduleTask'
})

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const EnableOrNotOptions = ref([])
const ScheduleTaskStatusOptions = ref([])
const formData = ref({
        taskDesc: '',
        asnNumber: '',
        asnDesc: '',
        crontabStr: '',
        taskConfig: '',
        enable: '0',
        taskStatus: '0',
        })
const configSet = `//请检查好任务配置再提交定时任务
//Json配置
{


}
`
const taskConfigTemp =  `{
        "scanDesc": "扫描AS906",
        "scanType": "1",
        "asnNumber": "AS906",
        "ipbatchSize": 100000,
        "enableTls": "1",
        "scanPorts": "443",
        "scanRate": 20000,
        "ipcheckThread": 100,
        "enableSpeedtest": "1",
}`


// 验证规则
const rule = reactive({
               taskDesc : [{
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
               asnNumber : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              },
               {
                 validator: (rule, value, callback) => {
                   const regex = /^AS\d+$/
                   if (!regex.test(value)) {
                     callback(new Error('ASN编号必须以 AS 开头，后跟纯数字'))
                   } else {
                     callback()
                   }
                 }
               },
              ],
               asnDesc : [{
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
               crontabStr : [{
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
               taskConfig : [{
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
// 排序
const sortChange = ({ prop, order }) => {
  const sortMap = {
            asnNumber: 'asn_number',
            ID:'id'
  }

  let sort = sortMap[prop]
  if(!sort){
   sort = prop.replace(/[A-Z]/g, match => `_${match.toLowerCase()}`)
  }

  searchInfo.value.sort = sort
  searchInfo.value.order = order
  getTableData()
}

const switchEnable = async(row) => {
  const currentRow = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    ...currentRow
  }
  const res = await updateScheduleTask(req)
  if (res.code === 0) {
    ElMessage({ type: 'success', message: `${req.enable === "0" ? '禁用' : '启用'}成功` })
    await getTableData()
  }
}

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
  const table = await getScheduleTaskList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    EnableOrNotOptions.value = await getDictFunc('EnableOrNot')
    ScheduleTaskStatusOptions.value = await getDictFunc('ScheduleTaskStatus')
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
            deleteScheduleTaskFunc(row)
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
      const res = await deleteScheduleTaskByIds({ IDs })
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
const updateScheduleTaskFunc = async(row) => {
    const res = await findScheduleTask({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteScheduleTaskFunc = async (row) => {
    const res = await deleteScheduleTask({ ID: row.ID })
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
        taskDesc: '',
        asnNumber: '',
        asnDesc: '',
        crontabStr: '',
        taskConfig: '',
        enable: '0',
        taskStatus: '0',
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  ElMessageBox.confirm('确定要创建此定时任务吗?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                  }).then(async() => {
                    res = await createScheduleTask(formData.value)
                    if (res.code === 0) {
                      ElMessage({
                        type: 'success',
                        message: '创建/更改成功'
                      })
                      closeDialog()
                      getTableData()
                    }
                  })
                  return
                case 'update':
                  res = await updateScheduleTask(formData.value)
                  break
                default:
                  res = await createScheduleTask(formData.value)
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

<style scoped>
.row-form {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
}

.row-form .el-form-item {
  width: 48%;
}

</style>
