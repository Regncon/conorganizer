import { BehaviorSubject } from 'rxjs';

const displayPool$ = new BehaviorSubject(false);
const childFriendly$ = new BehaviorSubject(false);
const possiblyEnglish$ = new BehaviorSubject(false);
const volunteersPossible$ = new BehaviorSubject(false);
const beginnerFriendly$ = new BehaviorSubject(false);
// Game types
const roleplaying$ = new BehaviorSubject(false);
const boardgame$ = new BehaviorSubject(false);
const other$ = new BehaviorSubject(false);
