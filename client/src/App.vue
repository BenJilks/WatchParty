<script setup lang="ts">
import Seats from "@/components/Seats.vue";
import Stage from "@/components/Stage.vue";
</script>

<template>
  <div class="background">
    <!--
    <video class="screen">
      <source src="/test.mp4" type="video/mp4">
    </video>
    !-->
    <img src="/test.png" class="screen" alt="" />
    <Stage />
    <Seats />
  </div>
</template>

<style scoped>
  @import 'assets/main.css';

  .background {
    position: fixed;
    width: 100%;
    height: 100%;
    left: 0;
    top: 0;

    background-color: #050505;
  }

  .screen {
    position: absolute;
    bottom: calc(var(--seat-height) + 7vh);

    object-fit: contain;
    width: calc(100% - 25vw);
    height: calc(100% - var(--seat-height) - 9vh);
    left: 50%;
    transform: translateX(-50%);

    background-color: #bbbf;
    box-shadow: 0 0 5vh 1vh #fff6;
    border-radius: 0.5vh;
  }
</style>
