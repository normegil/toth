<template>
  <header class="title-bar">
    <Search class="title-bar__search" />
    <Notification />
    <UserInfo class="title-bar__user-info"/>
  </header>
</template>

<script lang="ts">
  import Vue from "vue";
  import Component from "vue-class-component";
  import Search from "../../components/title-bar/Search.vue";
  import Notification from "../../components/navigation/Notification.vue";
  import UserInfo from "../../components/title-bar/UserInfo.vue";

@Component({
    components: {UserInfo, Notification, Search }
})
export default class TitleBar extends Vue {}
</script>

<style lang="scss">
.title-bar {
  background-color: $color-white;
  display: flex;
  border-bottom: $color-grey-7 solid 1px;
  color: $color-grey-1;

  &__search {
    flex-grow: 2;
    margin: 0.75rem;
  }

  &__notification {
    margin: 0.9rem;
  }

  &__user-info {
    margin: 0.9rem;
  }
}
</style>
