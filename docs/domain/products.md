## Product

--

A product is the basic unit we are able to put on the system, every other
relationship to what we can sell or display will start from here it can
be sales it be rate it can be offers.

The basic premise is that an user can create products under a workspace
and it can have category, prices, offers, sales attached to the product,
images, details, descriptions, videos, reviews etc so the product needs
to be the most basic unit and in a sense should be atomic

A product should always be under a workspace and not a user, why this is made
so is because different workspace can have duplicate products and the namespace
is made by workspace. Workspace can be considered as organizations and users
can be added with privileges and if they have product edit privilege they could
edit products/create privilege to edit products.

A product should contain the following a product name, product category,
product_id, sku_id, price, description, images and so much more. This causes a
huge problem that the product needs to be split into multiple tables and that
is a big problem so we will put that off for now and just build the bare
minimum with a single table it will contain the details that is needed in a
single table

The first step is to create a product, it looks like we should create a product
category first but I think it should be product because that is the core part
of the system and without that we can't move forward so creating product is the
priority rather than creating category and spending time there without anything
to push forward. We can always update the migrations and that is not a problem,
so building a product section is the first step for us at this stage.

The product table should contain the following properties:

1. sku id
2. product id
3. product name
4. product price
5. product description
6. workspace id
7. created at
8. edited at
9. deleted at
