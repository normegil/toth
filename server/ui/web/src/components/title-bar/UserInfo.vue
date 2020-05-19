<template>
  <a class="user-info-button">
    {{ user.name }}
    <CaretLeft class="user-info-button__icon" />
  </a>
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
  get user(): User {
    return this.$store.state.auth.authentifiedUser;
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

  &:hover {
    color: $color-complementary-light;
  }

  &:active {
    color: $color-complementary;
  }

  &__icon {
    margin-left: 0.5rem;
    position: relative;
    top: 0.3rem;
    height: 1.5rem;
  }
}
</style>
