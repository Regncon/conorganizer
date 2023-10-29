import { CustomEventTypeNames } from './enums';
import { CustomEventTypeFilteredEvent } from './types';

declare global {
    interface WindowEventHandlersEventMap {
        [CustomEventTypeNames.FilterChanges]: CustomEvent<CustomEventTypeFilteredEvent>;
    }
}
