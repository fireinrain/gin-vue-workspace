<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="ASN名称:" prop="asnName">
          <el-input v-model="formData.asnName" :clearable="true"  placeholder="请输入ASN名称" />
       </el-form-item>
        <el-form-item label="定时任务ID:" prop="scheduleTaskId">
          <el-input v-model.number="formData.scheduleTaskId" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="起始时间:" prop="startTime">
          <el-date-picker v-model="formData.startTime" type="date" placeholder="选择日期" :clearable="true"></el-date-picker>
       </el-form-item>
        <el-form-item label="结束时间:" prop="endTime">
          <el-date-picker v-model="formData.endTime" type="date" placeholder="选择日期" :clearable="true"></el-date-picker>
       </el-form-item>
        <el-form-item label="耗时:" prop="costTime">
          <el-input v-model.number="formData.costTime" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="任务状态:" prop="histStatus">
           <el-select v-model="formData.histStatus" placeholder="请选择任务状态" style="width:100%" :clearable="true" >
              <el-option v-for="(item,key) in ScheduleHistStatusOptions" :key="key" :label="item.label" :value="item.value" />
           </el-select>
       </el-form-item>
        <el-form-item label="任务结果:" prop="taskResult">
          <el-input v-model="formData.taskResult" :clearable="true"  placeholder="请输入任务结果" />
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
import {createScheduleTaskHist, findScheduleTaskHist, updateScheduleTaskHist} from '@/api/cfscan/scheduleTaskHist'
// 自动获取字典
import {getDictFunc} from '@/utils/format'
import {useRoute, useRouter} from "vue-router"
import {ElMessage} from 'element-plus'
import {reactive, ref} from 'vue'

defineOptions({
    name: 'ScheduleTaskHistForm'
})

const route = useRoute()
const router = useRouter()

const type = ref('')
const ScheduleHistStatusOptions = ref([])
const formData = ref({
            asnName: '',
            scheduleTaskId: undefined,
            startTime: new Date(),
            endTime: new Date(),
            costTime: undefined,
            histStatus: '',
            taskResult: '',
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findScheduleTaskHist({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    ScheduleHistStatusOptions.value = await getDictFunc('ScheduleHistStatus')
}

init()
// 保存按钮
const save = async() => {
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return
            let res
           switch (type.value) {
             case 'create':
               res = await createScheduleTaskHist(formData.value)
               break
             case 'update':
               res = await updateScheduleTaskHist(formData.value)
               break
             default:
               res = await createScheduleTaskHist(formData.value)
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
