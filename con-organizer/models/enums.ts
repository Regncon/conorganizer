export enum Pool {
    'none' = 'Ingen',
    'FridayEvening' = 'Fredag Kveld',
    'SaturdayMorning' = 'Lørdag Morgen',
    'SaturdayEvening' = 'Lørdag Kveld',
    'SundayMorning' = 'Søndag Morgen',
}

export enum GameType {
    'none' = 'Ingen',
    'roleplaying' = 'Rollespill',
    'boardgame' = 'Brettspill',
    'other' = 'Annet',
}
export enum FirebaseCollections {
    userSetting = 'usersettings',
    events = 'events',
    Test = 'Test',
    Participants = 'participants',
    Enrollments = 'enrollments',
    EventParticipants = 'eventParticipants',
    EnrollmentChoices = 'enrollmentChoices',
}

export enum EnrollmentOptions {
    'NotInterested',
    'IfIHaveTo',
    'Interested',
    'VeryInterested',
}
export enum GetLoginInfoResponse {
    Exists = 'Exists',
    Created = 'created',
    DontExist = "Don't exist",
}
