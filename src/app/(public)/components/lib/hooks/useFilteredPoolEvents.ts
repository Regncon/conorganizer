// useFilteredPoolEvents.ts
import useReadLocalStorage from '$lib/hooks/useReadLocalStorage';
import type { PoolEvent } from '$lib/types';
import { useMemo } from 'react';

// Define the structure of your filters
interface FilterCriteria {
    isActive: boolean;
    // You can extend this with more properties if needed
}

interface Filters {
    childFriendly: FilterCriteria;
    adultsOnly: FilterCriteria;
    beginnerFriendly: FilterCriteria;
    lessThanThreeHours: FilterCriteria;
    moreThanSixHours: FilterCriteria;
    possiblyEnglish: FilterCriteria;
    // Add more filters here if necessary
}

// Mapping between filter keys and PoolEvent properties
const filterMapping: { [key in keyof Filters]: keyof PoolEvent } = {
    childFriendly: 'childFriendly',
    adultsOnly: 'adultsOnly',
    beginnerFriendly: 'beginnerFriendly',
    lessThanThreeHours: 'lessThanThreeHours',
    moreThanSixHours: 'moreThanSixHours',
    possiblyEnglish: 'possiblyEnglish',
    // Extend mapping if you add more filters
};

// Function to remove unpublished events
const removeUnpublishedEventsFilter = (event: PoolEvent) => event.published;

// Refactored filteredEvents function with Logical OR
export const filteredEvents = (events: PoolEvent[], filters: Filters): PoolEvent[] => {
    // First, exclude unpublished events
    const publishedEvents = events.filter(removeUnpublishedEventsFilter);

    // Collect active filters
    const activeFilters = Object.entries(filters)
        .filter(([_, filter]) => filter.isActive)
        .map(([filterKey, _]) => filterKey as keyof Filters);

    // If no filters are active, return all published events
    if (activeFilters.length === 0) {
        return publishedEvents;
    }

    // Filter events that match any active filter
    return publishedEvents.filter((event) => {
        return activeFilters.some((filterKey) => {
            const eventProperty = event[filterMapping[filterKey]];
            return eventProperty;
        });
    });
};

// Hook to use filtered events with memoization
export const useFilteredPoolEvents = (events: PoolEvent[]) => {
    const filters = useReadLocalStorage('filters') as Filters;

    const filtered = useMemo(() => {
        // Debugging logs
        console.log('Applying Filters:', filters);
        console.log('Total Events Before Filtering:', events.length);
        console.log(
            'Total Events Before Filtering:',
            events.map((e) => ({
                childFriendly: e.childFriendly,
                adultsOnly: e.adultsOnly,
                beginnerFriendly: e.beginnerFriendly,
                lessThanThreeHours: e.lessThanThreeHours,
                moreThanSixHours: e.moreThanSixHours,
                possiblyEnglish: e.possiblyEnglish,
                title: e.title,
            }))
        );

        const result = filteredEvents(events, filters);

        console.log('Total Events After Filtering:', result.length);
        console.log('Filtered Events:', result);

        return result;
    }, [events, filters]);

    return filtered;
};
