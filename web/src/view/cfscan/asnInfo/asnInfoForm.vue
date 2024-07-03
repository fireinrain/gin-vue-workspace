<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="top" :inline="true" :rules="rule" label-width="80px">
        <el-form-item label="ASN名称:" prop="asnName" style="width: 30%">
          <el-input v-model="formData.asnName" :clearable="true"  placeholder="请输入ASN名称" />
       </el-form-item>
        <el-form-item label="ASN全名:" prop="fullName" style="width: 30%">
          <el-input v-model="formData.fullName" :clearable="true"  placeholder="请输入ASN全名" />
       </el-form-item>
        <el-form-item label="IPV4数量:" prop="ipv4Counts" style="width: 30%">
          <el-input v-model.number="formData.ipv4Counts" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="IPV6数量:" prop="ipv6Counts" style="width: 30%">
          <el-input v-model.number="formData.ipv6Counts" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="节点数量:" prop="peersCounts" style="width: 30%">
          <el-input v-model.number="formData.peersCounts" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="IPV4节点数量:" prop="ipv4Peers" style="width: 30%">
          <el-input v-model.number="formData.ipv4Peers" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="IPV6节点数量:" prop="ipv6Peers" style="width: 30%">
          <el-input v-model.number="formData.ipv6Peers" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="IP前缀数量:" prop="prefixesCounts" style="width: 30%">
          <el-input v-model.number="formData.prefixesCounts" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="IPV4前缀数量:" prop="ipv4Prefixies" style="width: 30%">
          <el-input v-model.number="formData.ipv4Prefixies" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="IPV6前缀数量:" prop="ipv6Prefixies" style="width: 30%">
          <el-input v-model.number="formData.ipv6Prefixies" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="地区登记:" prop="regionalRegistry" style="width: 30%">
          <el-input v-model="formData.regionalRegistry" :clearable="true"  placeholder="请输入地区登记" />
       </el-form-item>
        <el-form-item label="带宽估算:" prop="trafficBandwidth" style="width: 30%">
          <el-input v-model="formData.trafficBandwidth" :clearable="true"  placeholder="请输入带宽估算" />
       </el-form-item>
        <el-form-item label="分配状态:" prop="allocationStatus" style="width: 30%">
          <el-input v-model="formData.allocationStatus" :clearable="true"  placeholder="请输入分配状态" />
       </el-form-item>
        <el-form-item label="流量比率:" prop="trafficRatio" style="width: 30%">
          <el-input v-model="formData.trafficRatio" :clearable="true"  placeholder="请输入流量比率" />
       </el-form-item>
        <el-form-item label="分配日期:" prop="allocationDate" style="width: 30%">
          <el-input v-model="formData.allocationDate" :clearable="true"  placeholder="请输入分配日期" />
       </el-form-item>
        <el-form-item label="官方网址:" prop="website" style="width: 30%">
          <el-input v-model="formData.website" :clearable="true"  placeholder="请输入官方网址" />
       </el-form-item>
        <el-form-item label="分配国家:" prop="allocationCountry" style="width: 30%">
          <el-input v-model="formData.allocationCountry" :clearable="true"  placeholder="请输入分配国家" />
       </el-form-item>
        <el-form-item label="IPV4 CIDR:" prop="ipv4CIDR" style="width: 30%">
          <el-input v-model="formData.ipv4CIDR" :clearable="true"  placeholder="请输入IPV4 CIDR" />
       </el-form-item>
        <el-form-item label="是否开启:" prop="enable" style="width: 30%">
          <el-input v-model.number="formData.enable" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="CIDR最后更新时间:" prop="lastCIDRUpdate" style="width: 30%">
          <el-date-picker v-model="formData.lastCIDRUpdate" type="date" placeholder="选择日期" :clearable="true"></el-date-picker>
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
import {createAsnInfo, findAsnInfo, updateAsnInfo} from '@/api/cfscan/asnInfo'
// 自动获取字典
import {useRoute, useRouter} from "vue-router"
import {ElMessage} from 'element-plus'
import {reactive, ref} from 'vue'

defineOptions({
    name: 'AsnInfoForm'
})

const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
            asnName: '',
            fullName: '',
            ipv4Counts: undefined,
            ipv6Counts: undefined,
            peersCounts: undefined,
            ipv4Peers: undefined,
            ipv6Peers: undefined,
            prefixesCounts: undefined,
            ipv4Prefixies: undefined,
            ipv6Prefixies: undefined,
            regionalRegistry: '',
            trafficBandwidth: '',
            allocationStatus: '',
            trafficRatio: '',
            allocationDate: '',
            website: '',
            allocationCountry: '',
            ipv4CIDR: '',
            enable: undefined,
            lastCIDRUpdate: new Date(),
        })
// 验证规则
const rule = reactive({
               asnName : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               enable : [{
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
      const res = await findAsnInfo({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return
            let res
           switch (type.value) {
             case 'create':
               res = await createAsnInfo(formData.value)
               break
             case 'update':
               res = await updateAsnInfo(formData.value)
               break
             default:
               res = await createAsnInfo(formData.value)
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
