import type { IconName, PoolEvent } from '$lib/types';
import { useEffect, useState } from 'react';

const removeUnpublishedEventsFilter = (event: PoolEvent) => event.published;

// export const useFilteredEvents = (events: PoolEvent[], searchParams: { [key in IconName]: string }) => {
//     useEffect(() => {
//         console.log('here');
//     }, [filteredEvents]);
//     console.log(filteredEvents(events, searchParams), 'filtered events');

//     return filteredEvents(events, searchParams);
// };
export const filteredEvents = (events: PoolEvent[], searchParams: { [key in IconName]: string }): PoolEvent[] => {
    console.log(
        events.filter(removeUnpublishedEventsFilter),
        // .filter((event) => childFriendlyFilter(event, searchParams))
        // .filter((event) => adultsOnlyFilter(event, searchParams))
        // .filter((event) => beginnerFriendlyFilter(event, searchParams))
        // .filter((event) => lessThanThreeHoursFilter(event, searchParams))
        // .filter((event) => moreThanSixHoursFilter(event, searchParams))
        // .filter((event) => possiblyEnglishFilter(event, searchParams)),
        'filtered events'
    );

    return events.filter(removeUnpublishedEventsFilter).filter((event) => childFriendlyFilter(event, searchParams));
    // .filter((event) => adultsOnlyFilter(event, searchParams))
    // .filter((event) => beginnerFriendlyFilter(event, searchParams))
    // .filter((event) => lessThanThreeHoursFilter(event, searchParams))
    // .filter((event) => moreThanSixHoursFilter(event, searchParams))
    // .filter((event) => possiblyEnglishFilter(event, searchParams));
};
export const childFriendlyFilter = (event: PoolEvent, searchParams: { [key in IconName]: string }) => {
    return searchParams.childFriendly === 'true' ? true : false;
};
export const adultsOnlyFilter = (event: PoolEvent, searchParams: { [key in IconName]: string }) => {
    return searchParams.adultsOnly === 'true' ? true : false;
};
export const beginnerFriendlyFilter = (event: PoolEvent, searchParams: { [key in IconName]: string }) => {
    return searchParams.beginnerFriendly === 'true' ? true : false;
};
export const lessThanThreeHoursFilter = (event: PoolEvent, searchParams: { [key in IconName]: string }) => {
    return searchParams.lessThanThreeHours === 'true' ? true : false;
};
export const moreThanSixHoursFilter = (event: PoolEvent, searchParams: { [key in IconName]: string }) => {
    return searchParams.moreThanSixHours === 'true' ? true : false;
};
export const possiblyEnglishFilter = (event: PoolEvent, searchParams: { [key in IconName]: string }) => {
    return searchParams.possiblyEnglish === 'true' ? true : false;
};
