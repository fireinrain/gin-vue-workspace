<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="scanDesc字段:" prop="scanDesc">
          <el-input v-model="formData.scanDesc" :clearable="true"  placeholder="请输入scanDesc字段" />
       </el-form-item>
        <el-form-item label="scanType字段:" prop="scanType">
           <el-select v-model="formData.scanType" placeholder="请选择scanType字段" style="width:100%" :clearable="true" >
              <el-option v-for="(item,key) in CFScanTypeOptions" :key="key" :label="item.label" :value="item.value" />
           </el-select>
       </el-form-item>
        <el-form-item label="asnNumber字段:" prop="asnNumber">
          <el-input v-model="formData.asnNumber" :clearable="true"  placeholder="请输入asnNumber字段" />
       </el-form-item>
        <el-form-item label="ipinfoType字段:" prop="ipinfoType">
           <el-select v-model="formData.ipinfoType" placeholder="请选择ipinfoType字段" style="width:100%" :clearable="true" >
              <el-option v-for="(item,key) in IPInfoTypeOptions" :key="key" :label="item.label" :value="item.value" />
           </el-select>
       </el-form-item>
        <el-form-item label="ipinfoList字段:" prop="ipinfoList">
          <el-input v-model="formData.ipinfoList" :clearable="true"  placeholder="请输入ipinfoList字段" />
       </el-form-item>
        <el-form-item label="enableTls字段:" prop="enableTls">
           <el-select v-model="formData.enableTls" placeholder="请选择enableTls字段" style="width:100%" :clearable="true" >
              <el-option v-for="(item,key) in TLSScanTypeOptions" :key="key" :label="item.label" :value="item.value" />
           </el-select>
       </el-form-item>
        <el-form-item label="scanPorts字段:" prop="scanPorts">
          <el-input v-model="formData.scanPorts" :clearable="true"  placeholder="请输入scanPorts字段" />
       </el-form-item>
        <el-form-item label="scanRate字段:" prop="scanRate">
          <el-input v-model.number="formData.scanRate" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="ipcheckThread字段:" prop="ipcheckThread">
          <el-input v-model.number="formData.ipcheckThread" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="enableSpeedtest字段:" prop="enableSpeedtest">
           <el-select v-model="formData.enableSpeedtest" placeholder="请选择enableSpeedtest字段" style="width:100%" :clearable="true" >
              <el-option v-for="(item,key) in IPSpeedTestOptions" :key="key" :label="item.label" :value="item.value" />
           </el-select>
       </el-form-item>
        <el-form-item label="scanStatus字段:" prop="scanStatus">
           <el-select v-model="formData.scanStatus" placeholder="请选择scanStatus字段" style="width:100%" :clearable="true" >
              <el-option v-for="(item,key) in ScanTaskStatusOptions" :key="key" :label="item.label" :value="item.value" />
           </el-select>
       </el-form-item>
        <el-form-item label="scanResult字段:" prop="scanResult">
          <el-input v-model="formData.scanResult" :clearable="true"  placeholder="请输入scanResult字段" />
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
import {createSubmitScan, findSubmitScan, updateSubmitScan} from '@/api/cfscan/submitScan'
// 自动获取字典
import {getDictFunc} from '@/utils/format'
import {useRoute, useRouter} from "vue-router"
import {ElMessage} from 'element-plus'
import {reactive, ref} from 'vue'

defineOptions({
    name: 'SubmitScanForm'
})

const route = useRoute()
const router = useRouter()

const type = ref('')
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
            enableTls: '',
            scanPorts: '',
            scanRate: undefined,
            ipcheckThread: undefined,
            enableSpeedtest: '',
            scanStatus: '',
            scanResult: '',
        })
// 验证规则
const rule = reactive({
               scanDesc : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               scanType : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               enableTls : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               scanPorts : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               scanRate : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               ipcheckThread : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               enableSpeedtest : [{
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
      const res = await findSubmitScan({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    TLSScanTypeOptions.value = await getDictFunc('TLSScanType')
    IPSpeedTestOptions.value = await getDictFunc('IPSpeedTest')
    ScanTaskStatusOptions.value = await getDictFunc('ScanTaskStatus')
    CFScanTypeOptions.value = await getDictFunc('CFScanType')
    IPInfoTypeOptions.value = await getDictFunc('IPInfoType')
}

init()
// 保存按钮
const save = async() => {
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
