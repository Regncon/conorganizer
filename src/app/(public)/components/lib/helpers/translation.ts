import { PoolName } from '$lib/enums';

export const translatedDays = new Map<PoolName, string>([
    [PoolName.fridayEvening, 'Fredag kveld'],
    [PoolName.saturdayMorning, 'Lørdag morgen'],
    [PoolName.saturdayEvening, 'Lørdag kveld'],
    [PoolName.sundayMorning, 'Søndag morgen'],
]);

export const getTranslatedDay = (day: PoolName) => translatedDays.get(day);
