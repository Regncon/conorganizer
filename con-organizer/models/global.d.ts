import { CustomEventTypeNames } from './enums';
import { CustomEventTypeFilteredEvent } from './types';

declare global {
    interface WindowEventMap {
        [CustomEventTypeNames.FilterChanges]: CustomEvent<CustomEventTypeFilteredEvent>;
    }
}
