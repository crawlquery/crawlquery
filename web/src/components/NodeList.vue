<script lang="ts" setup>
import { ref } from 'vue'
import axios from 'axios'
import { useAuthStore } from '../stores/auth'
import { Node } from '../types'
defineProps<{ nodes: Array<Node> }>()
const err = ref('')
</script>
<template>

    <div class="overflow-x-auto">
        <table class="table table-zebra">
            <!-- head -->
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Hostname</th>
                    <th>Port</th>
                    <th>Shard ID</th>
                    <th>Created At</th>
                </tr>
            </thead>
            <tbody>
                <!-- row 1 -->
                <tr v-for="node in nodes">
                    <td>{{ node.id }}</td>
                    <td>{{ node.hostname }}</td>
                    <td>{{ node.port }}</td>
                    <td>{{ node.shard_id }}</td>
                    <td>{{ node.created_at }}</td>
                </tr>

                <tr v-if="nodes.length === 0">
                    <td colspan="5" class="text-center">No nodes found</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>