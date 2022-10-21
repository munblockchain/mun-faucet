<template lang="pug">
.section
  .section-header
    h1 Community Links
  .section-main
    list-item(
      href="https://discord.gg/mun"
      target="_blank"
      title="MUN Community Chat"
      subtitle="#cosmos")
    list-item(
      href="https://blockexplorer.mun.money"
      target="_blank"
      title="MUN Testnet Explorer"
      subtitle="explorecosmos.network")
    list-item(
      href="https://mun.money"
      target="_blank"
      title="MUN Website"
      subtitle="cosmos.network")
</template>

<script>
import ListItem from "../components/NiListItem.vue";
export default {
  name: "faucet",
  components: {
    ListItem
  },
  data: () => ({}),
  methods: {}
};
</script>
