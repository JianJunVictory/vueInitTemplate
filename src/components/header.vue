<template>
    <div>
        <div class="announcement header__announcement">FREE SHIPPING ON NZ ORDERS OVER<span>$100.00</span>/ AU ORDERS OVER<span>$250.00</span></div>
        <div class="HDmenu">
           <ul>
            <li><a>MENS</a></li>
            <li><a>WOMENS</a></li>
            <li><a herf="javascript:;">FOOTWEAR</a></li>
            <li><a herf="javascript:;">BRANDS</a></li>
            <li><a herf="javascript:;">STORES</a></li>
            <li><a herf="javascript:;">CONTACT</a></li>
           </ul>
        </div>
    </div>
</template>
<script>
export default {
  name: 'Header',
  data () {
    return {
      centerDialogVisible: false
    }
  },
  methods: {
    overShow () {
      this.centerDialogVisible = true
    },
    noShow () {
      this.centerDialogVisible = false
    }
  }
}
</script>
<style scoped>
.header__announcement {
    grid-area: announce;
    justify-self: stretch;
}
@media (min-width: 1200px){
    .announcement {
        font-size: 1.1rem;
        letter-spacing: 1px;
    }
}
.announcement {
    background-color: black;
    color: white;
    font-size: 0.8rem;
    letter-spacing: .5px;
    line-height: 1.2;
    padding: 8px 10px;
    text-align: center;
}
/* 所有class为menu的div中的ul样式 */
div.HDmenu{
    margin-left:30%;
}
div.HDmenu ul
{
    list-style:none; /* 去掉ul前面的符号 */
    margin: 20px; /* 与外界元素的距离为0 */
    padding: 0px; /* 与内部元素的距离为0 */
    width: auto; /* 宽度根据元素内容调整 */
    position: relative;
}
/* 所有class为menu的div中的ul中的li样式 */
div.HDmenu ul li
{
    float:left; /* 向左漂移，将竖排变为横排 */
    padding-left:20px;
    font-size: 0.9rem;
    font-weight: bold;
}
div.HDmenu ul li a:hover
{
    text-decoration: none; /* 不显示超链接下划线 */
    cursor:pointer
}
div.HDmenu ul li:hover .subHDmenu { display: block; }
div.HDmenu ul li .subHDmenu { display: none; }
</style>
