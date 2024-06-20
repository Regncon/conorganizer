import EventCardBig from '$app/(public)/EventCardBig';
import Link from 'next/link';

const Events = async () => {
    return (
        <Link href="/admin/dashboard/events/event-dashboard/1">
            <EventCardBig
                gameMaster="Kåre Carlsson"
                shortDescription="is water wet? find out in this session! (18+) (NSFW) anyways do NOT join if you are under 18, this is a serious session for serious people only."
                system="DnD 4e"
                title="Kåres waterboarding session of doom and despair (18+)"
            />
        </Link>
    );
};

export default Events;
