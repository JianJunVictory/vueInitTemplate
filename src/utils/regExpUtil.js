function isPhone (phone) {
  /**
  * 中国电信号码格式验证 手机段： 133,153,180,181,189,177,1700,173
  * **/
  let CHINA_TELECOM_PATTERN = /(?:^(?:\+86)?1(?:33|53|7[37]|8[019])\d{8}$)|(?:^(?:\+86)?1700\d{7}$)/
  /**
  * 中国联通号码格式验证 手机段：130,131,132,155,156,185,186,145,176,1707,1708,1709,175
  * **/
  let CHINA_UNICOM_PATTERN = /(?:^(?:\+86)?1(?:3[0-2]|4[5]|5[56]|7[56]|8[56])\d{8}$)|(?:^(?:\+86)?170[7-9]\d{7}$)/
  /**
  * 简单手机号码校验，校验手机号码的长度和1开头
  */
  let SIMPLE_PHONE_CHECK = /^(?:\\+86)?1\d{10}$/
  /**
  * 中国移动号码格式验证
  * 手机段：134,135,136,137,138,139,150,151,152,157,158,159,182,183,184
  * ,187,188,147,178,1705
  *
  **/
  let CHINA_MOBILE_PATTERN = /(?:^(?:\+86)?1(?:3[4-9]|4[7]|5[0-27-9]|7[8]|8[2-478])\d{8}$)|(?:^(?:\+86)?1705\d{7}$)/
  if (SIMPLE_PHONE_CHECK.test(phone)) {
    if (CHINA_TELECOM_PATTERN.test(phone) || CHINA_MOBILE_PATTERN.test(phone) || CHINA_UNICOM_PATTERN.test(phone)) {
      return true
    } else {
      return false
    }
  } else {
    return false
  }
}

function isEmail (email) {
  let EMAIL_PATTERN = /^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$/
  if (EMAIL_PATTERN.test(email)) {
    return true
  }
  return false
}

function checkPassword (str) {
  let reg3 = /[a-zA-Z0-9]{8,16}/
  if (reg3.test(str)) {
    return true
  } else if (!reg3.test(str)) {
    // alert("需含有字母")
    return false
  }
}

export default {
  isPhone: isPhone,
  isEmail: isEmail,
  checkPassword: checkPassword
}
