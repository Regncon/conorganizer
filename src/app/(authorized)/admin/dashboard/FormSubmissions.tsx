import CardBase from '$app/(authorized)/dashboard/CardBase';

const FormSubmissions = async () => {
    return (
        <CardBase
            href="/admin/dashboard/form-submissions"
            subTitle="Trykk for å gå til innsendte skjema"
            img="/form-submissions-thumbnail.webp"
            imgAlt="Innsendte skjema"
            title="Innsendte skjema"
        />
    );
};

export default FormSubmissions;
