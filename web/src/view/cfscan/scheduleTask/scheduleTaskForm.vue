<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="任务描述:" prop="taskDesc">
          <el-input v-model="formData.taskDesc" :clearable="true"  placeholder="请输入任务描述" />
       </el-form-item>
        <el-form-item label="ASN编号:" prop="asnNumber">
          <el-input v-model="formData.asnNumber" :clearable="true"  placeholder="请输入ASN编号" />
       </el-form-item>
        <el-form-item label="ASN描述:" prop="asnDesc">
          <el-input v-model="formData.asnDesc" :clearable="true"  placeholder="请输入ASN描述" />
       </el-form-item>
        <el-form-item label="定时表达式:" prop="crontabStr">
          <el-input v-model="formData.crontabStr" :clearable="true"  placeholder="请输入定时表达式" />
       </el-form-item>
        <el-form-item label="任务配置:" prop="taskConfig">
          <el-input v-model="formData.taskConfig" :clearable="true"  placeholder="请输入任务配置" />
       </el-form-item>
        <el-form-item label="是否开启:" prop="enable">
           <el-select v-model="formData.enable" placeholder="请选择是否开启" style="width:100%" :clearable="true" >
              <el-option v-for="(item,key) in EnableOrNotOptions" :key="key" :label="item.label" :value="item.value" />
           </el-select>
       </el-form-item>
        <el-form-item label="任务状态:" prop="taskStatus">
           <el-select v-model="formData.taskStatus" placeholder="请选择任务状态" style="width:100%" :clearable="true" >
              <el-option v-for="(item,key) in ScheduleTaskStatusOptions" :key="key" :label="item.label" :value="item.value" />
           </el-select>
       </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {createScheduleTask, findScheduleTask, updateScheduleTask} from '@/api/cfscan/scheduleTask'
// 自动获取字典
import {getDictFunc} from '@/utils/format'
import {useRoute, useRouter} from "vue-router"
import {ElMessage} from 'element-plus'
import {reactive, ref} from 'vue'

defineOptions({
    name: 'ScheduleTaskForm'
})

const route = useRoute()
const router = useRouter()

const type = ref('')
const EnableOrNotOptions = ref([])
const ScheduleTaskStatusOptions = ref([])
const formData = ref({
            taskDesc: '',
            asnNumber: '',
            asnDesc: '',
            crontabStr: '',
            taskConfig: '',
            enable: '',
            taskStatus: '',
        })
// 验证规则
const rule = reactive({
               taskDesc : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               asnNumber : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               asnDesc : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               crontabStr : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               taskConfig : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findScheduleTask({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    EnableOrNotOptions.value = await getDictFunc('EnableOrNot')
    ScheduleTaskStatusOptions.value = await getDictFunc('ScheduleTaskStatus')
}

init()
// 保存按钮
const save = async() => {
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return
            let res
           switch (type.value) {
             case 'create':
               res = await createScheduleTask(formData.value)
               break
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
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
