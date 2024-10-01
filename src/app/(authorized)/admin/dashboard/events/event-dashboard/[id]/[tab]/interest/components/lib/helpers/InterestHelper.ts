import { InterestLevel } from '$lib/enums';
import AwakeDragons from 'public/interessedragene/2024AwakeDragons1_1.png';
import HappyDragons from 'public/interessedragene/2024HappyDragons1_1.png';
import SleepyDragons from 'public/interessedragene/2024SleepyDragons1_1.png';
import VeryHappyDragons from 'public/interessedragene/2024VeryHappyDragons1_1.png';

export const InterestLevelToLabel = {
    [InterestLevel.VeryInterested]: 'Veldig interessert',
    [InterestLevel.Interested]: 'Interessert',
    [InterestLevel.SomewhatInterested]: 'Litt interessert',
    [InterestLevel.NotInterested]: 'Ikke interessert',
};

export const interestLevelToImage = {
    [InterestLevel.VeryInterested]: VeryHappyDragons,
    [InterestLevel.Interested]: HappyDragons,
    [InterestLevel.SomewhatInterested]: AwakeDragons,
    [InterestLevel.NotInterested]: SleepyDragons,
};
