import { PoolName } from '$lib/enums';
import type { PoolPlayer } from '$lib/types';

export const generatePlayerInPoolMap = (playersInPool: PoolPlayer[]) => {
    const playerInPoolMap = new Map<PoolName, PoolPlayer>([
        [PoolName.fridayEvening, {} as PoolPlayer],
        [PoolName.saturdayMorning, {} as PoolPlayer],
        [PoolName.saturdayEvening, {} as PoolPlayer],
        [PoolName.sundayMorning, {} as PoolPlayer],
    ]);
    playersInPool.forEach((player) => {
        playerInPoolMap.set(player.poolName, player);
    });
    return playerInPoolMap;
};
