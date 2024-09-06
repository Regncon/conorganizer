import CardBase from '$app/(authorized)/dashboard/components/CardBase';

const Events = async () => {
    return (
        <CardBase
            href="/admin/dashboard/events"
            subTitle="Trykk for å gå til alle arrangement"
            img="/my-events.jpg"
            imgAlt="Alle arrangement"
            title="Alle arrangement"
        />
    );
};

export default Events;
