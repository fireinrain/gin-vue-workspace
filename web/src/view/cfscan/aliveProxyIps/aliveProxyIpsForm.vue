<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="ASN编号:" prop="asnNumber">
          <el-input v-model="formData.asnNumber" :clearable="true"  placeholder="请输入ASN编号" />
       </el-form-item>
        <el-form-item label="ASN描述:" prop="asnDesc">
          <el-input v-model="formData.asnDesc" :clearable="true"  placeholder="请输入ASN描述" />
       </el-form-item>
        <el-form-item label="IP地址:" prop="ip">
          <el-input v-model="formData.ip" :clearable="true"  placeholder="请输入IP地址" />
       </el-form-item>
        <el-form-item label="端口:" prop="port">
          <el-input v-model.number="formData.port" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="开启TLS:" prop="enableTls">
           <el-select v-model="formData.enableTls" placeholder="请选择开启TLS" style="width:100%" :clearable="true" >
              <el-option v-for="(item,key) in TLSDisplayEnableOptions" :key="key" :label="item.label" :value="item.value" />
           </el-select>
       </el-form-item>
        <el-form-item label="物理距离:" prop="geoDistance">
          <el-input v-model.number="formData.geoDistance" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="数据中心:" prop="dataCenter">
          <el-input v-model="formData.dataCenter" :clearable="true"  placeholder="请输入数据中心" />
       </el-form-item>
        <el-form-item label="地区:" prop="region">
          <el-input v-model="formData.region" :clearable="true"  placeholder="请输入地区" />
       </el-form-item>
        <el-form-item label="城市:" prop="city">
          <el-input v-model="formData.city" :clearable="true"  placeholder="请输入城市" />
       </el-form-item>
        <el-form-item label="延迟:" prop="latency">
          <el-input v-model="formData.latency" :clearable="true"  placeholder="请输入延迟" />
       </el-form-item>
        <el-form-item label="TCP延迟:" prop="tcpDuration">
          <el-input v-model.number="formData.tcpDuration" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="下载速度:" prop="downloadSpeed">
          <el-input v-model="formData.downloadSpeed" :clearable="true"  placeholder="请输入下载速度" />
       </el-form-item>
        <el-form-item label="存活时间:" prop="ttl">
          <el-input v-model.number="formData.ttl" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="IP描述:" prop="descStr">
          <el-input v-model="formData.descStr" :clearable="true"  placeholder="请输入IP描述" />
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
import {createAliveProxyIps, findAliveProxyIps, updateAliveProxyIps} from '@/api/cfscan/aliveProxyIps'
// 自动获取字典
import {getDictFunc} from '@/utils/format'
import {useRoute, useRouter} from "vue-router"
import {ElMessage} from 'element-plus'
import {reactive, ref} from 'vue'

defineOptions({
    name: 'AliveProxyIpsForm'
})

const route = useRoute()
const router = useRouter()

const type = ref('')
const TLSDisplayEnableOptions = ref([])
const formData = ref({
            asnNumber: '',
            asnDesc: '',
            ip: '',
            port: undefined,
            enableTls: '',
            geoDistance: undefined,
            dataCenter: '',
            region: '',
            city: '',
            latency: '',
            tcpDuration: undefined,
            downloadSpeed: '',
            ttl: undefined,
            descStr: '',
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findAliveProxyIps({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    TLSDisplayEnableOptions.value = await getDictFunc('TLSDisplayEnable')
}

init()
// 保存按钮
const save = async() => {
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return
            let res
           switch (type.value) {
             case 'create':
               res = await createAliveProxyIps(formData.value)
               break
             case 'update':
               res = await updateAliveProxyIps(formData.value)
               break
             default:
               res = await createAliveProxyIps(formData.value)
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
