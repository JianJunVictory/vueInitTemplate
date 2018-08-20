<template>
    <div>
    {{info}}
        <el-form :model="loginFrom" status-icon :rules="rules2" ref="loginFrom" label-width="100px" class="demo-ruleForm">
        <el-form-item :label="$t('login.email')" prop="email">
            <el-input type="email" v-model="loginFrom.email" auto-complete="off" :placeholder="$t('login.email')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('login.password')" prop="password">
            <el-input type="password" v-model="loginFrom.password" auto-complete="off" :placeholder="$t('login.password')"></el-input>
        </el-form-item>
        <el-form-item>
            <el-button type="primary" @click="submitForm('loginFrom')">{{$t('login.submit')}}</el-button>
            <el-button @click="resetForm('loginFrom')">{{$t('login.reset')}}</el-button>
        </el-form-item>
        </el-form>
    </div>
</template>

<script>
import { mapState } from 'vuex'
import regExpUtil from '../utils/regExpUtil'
export default {
  name: 'Login',
  data () {
    var validateEmail = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入手机号'))
      } else if (!regExpUtil.isEmail(value)) {
        callback(new Error('请输入正确手机号'))
      } else {
        callback()
      }
    }
    var validatePassword = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入密码'))
      } else if (!regExpUtil.checkPassword(value)) {
        callback(new Error('请输入至少由数字,字母组成的8到16位密码'))
      } else {
        callback()
      }
    }
    return {
      title: this.$store.state.info,
      loginFrom: {
        email: '',
        password: ''
      },
      rules2: {
        email: [
          { validator: validateEmail, trigger: 'blur' }
        ],
        password: [
          { validator: validatePassword, trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    submitForm (formName) {
    },
    resetForm (formName) {
    }
  },
  computed: {
    ...mapState([
      'info'
    ])
  }
}
</script>

<style lang="">
</style>
