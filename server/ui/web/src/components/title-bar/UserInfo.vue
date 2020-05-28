<template>
  <div class="user-info-button">
    {{ user.name }}
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import User from "../../model/User";
import CaretLeft from "../../assets/images/icons/caret-left-solid.svg";

@Component({
  components: { CaretLeft }
})
export default class UserInfo extends Vue {
  get user(): User | undefined {
    return this.$store.state.auth.authentifiedUser;
  }

  get username(): string {
    let user = this.user;
    if (user === undefined) {
      return ""
    }
    return user.name
  }
}
</script>

<style lang="scss">
.user-info-button {
  display: block;
  font-weight: bold;
  font-size: 1.2rem;
  padding: 1.5rem;
  transition: all 0.2s;
  cursor: pointer;
}
</style>
