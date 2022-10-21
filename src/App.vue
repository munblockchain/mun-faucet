<template lang="pug">
#app
  tm-cookie-consent
  .top-bar
    .left
      a(href="#") #[img(src="~assets/icon.svg" width="45" height="45")]
      .value {{ config.chain }}
    .right
      a(href="https://blockexplorer.mun.money/" target="_blank") explorer #[i.material-icons assessment]
    .cube
      img(src="~assets/cube.png" width="250" height="250")
  router-view
  notifications(:notifications='notifications' theme='cosmos')
  #bottom &copy; 2022 MUN Blockchain 
</template>


<script>
import { mapGetters } from "vuex";
import { TmCookieConsent } from "@tendermint/ui";
import Notifications from "@nylira/vue-notifications";

export default {
  name: "app",
  components: {
    Notifications,
    TmCookieConsent
  },
  computed: {
    ...mapGetters(["notifications", "config"])
  },
  mounted() {}
};
</script>

<style lang="stylus" src="./styles/app.styl"></style>

<style lang="stylus">
@import '~variables'


.top-bar
  position fixed
  top 0
  right 0
  background var(--accent)
  text-align center
  height 3rem
  width 100%
  z-index 100

  display flex
  justify-content space-between
  line-height 3rem

  .left
    display flex
    padding 0 1rem
  .value
    font-weight bold
    color var(--bright)
  .right
    a
      display block
      background var(--mc)
      padding 0 1rem
      color var(--bright)
      display flex
      align-items center
      i
        margin-left 0.5rem
  .cube
    position absolute
    img
      position absolute
      top 4rem
      left 4rem
      z-index 1
#bottom
  color var(--bc)
  font-size 0.75rem
  padding 1rem 0
  text-align: center;
</style>
