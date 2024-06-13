import CardBase from '$app/(authorized)/dashboard/CardBase';

const Events = async () => {
    return (
        <CardBase
            href="/admin/dashboard/events"
            subTitle="Trykk for å gå til alle arrangemanter"
            img="/my-events.jpg"
            imgAlt="Alle arrangemanter"
            title="Alle arrangemanter"
        />
    );
};

export default Events;
