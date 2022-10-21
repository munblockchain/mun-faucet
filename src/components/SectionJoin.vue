<template lang="pug">
.section
  .section-header
    h1 Don't have a testnet address?
    p Connect wallet to MUN to have mun address.
  .section-main
    btn(type="anchor" href="http://wallet.mun.money/" value="Join the latest testnet" size="lg" target="_blank" icon="wifi_tethering")
</template>

<script>
import Btn from "@nylira/vue-button";
export default {
  name: "section-join",
  components: {
    Btn
  }
};
</script>
