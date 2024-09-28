import { PoolName } from '$lib/enums';

export const translatedDays = new Map<PoolName, string>([
    [PoolName.fridayEvening, 'Fredag kveld'],
    [PoolName.saturdayMorning, 'Lørdag morgen'],
    [PoolName.saturdayEvening, 'Lørdag kveld'],
    [PoolName.sundayMorning, 'Søndag morgen'],
]);

export const getTranslatedDay = (day: PoolName) => translatedDays.get(day);

export const translatedDaysAndTime = new Map<PoolName, string>([
    [PoolName.fridayEvening, 'Fredag kveld Kl 18 - 23'],
    [PoolName.saturdayMorning, 'Lørdag morgen Kl 10 - 15'],
    [PoolName.saturdayEvening, 'Lørdag kveld Kl 18 - 23'],
    [PoolName.sundayMorning, 'Søndag morgen Kl 10 - 15'],
]);

export const getTranslatedDayAndTime = (day: PoolName) => translatedDaysAndTime.get(day);
