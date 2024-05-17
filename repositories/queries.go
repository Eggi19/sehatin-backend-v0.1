package repositories

const (
	qCreateOneUser = `
		INSERT INTO users (name, email, birth_date, password)
		VALUES (--values--)
		RETURNING id;
	`
	qCreateUserVerificationToken = `
		INSERT INTO user_verification_tokens (token, expired_at, user_id)
		VALUES ($1, $2, $3);
	`

	qFindUserByEmail = `
		SELECT id, name, email, password, is_verified, is_google, is_online, profile_picture FROM users
		WHERE email = $1 AND deleted_at IS NULL;
	`

	qVerifyUser = `
		UPDATE users
		SET is_verified = true
		WHERE email = $1 AND deleted_at IS NULL;
	`

	qFindOneUserById = `
		SELECT u.id, u.name, u.email, u.birth_date, g.id ,g.name, u.is_verified, u.is_google, u.is_online, u.profile_picture
		FROM users u
		LEFT JOIN genders g ON u.gender_id = g.id
		WHERE u.id =$1 AND u.deleted_at IS NULL
		LIMIT 1;
	`

	qDeleteUserById = `
		UPDATE users SET
		deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	qUserColl = `
		, u.id, u.name, u.email, u.birth_date, g.id ,g.name, u.is_verified, u.is_google, u.is_online, u.profile_picture
	`

	qUserCommand = `
		FROM users u
		LEFT JOIN genders g ON u.gender_id = g.id
		WHERE u.deleted_at IS NULL
	`

	qUpdateOneUser = `
		UPDATE users SET name = $2, birth_date = $3, gender_id = $4, profile_picture = $5, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qUpdateOneUserWithOutProfile = `
		UPDATE users SET name = $2, birth_date = $3, gender_id = $4, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qUpadateUserPassword = `
		UPDATE users SET password = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qFindUserPasswordById = `
		SELECT id, password FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
)

const (
	qFindDoctorByEmail = `
		SELECT id, name, email, password, is_verified, is_google, is_online ,profile_picture FROM doctors
		WHERE email = $1 AND deleted_at IS NULL;
	`
	qFindDoctorById = `
		SELECT d.id, d.name, d.email, d.certificate, d.is_online, d.is_verified, d.is_google, d.fee, d.work_start_year, ds.id, ds.name, d.profile_picture FROM doctors d
		LEFT JOIN doctor_specialists ds ON d.doctor_specialists_id = ds.id
		WHERE d.id = $1 AND d.deleted_at IS NULL
		LIMIT 1;
	`

	qVerifyDoctor = `
		UPDATE doctors
		SET is_verified = true
		WHERE email = $1 AND deleted_at IS NULL;
	`
	qCreateOneDoctor = `
		INSERT INTO doctors (name, email, password, fee, certificate, work_start_year, doctor_specialists_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`
	qCreateDoctorVerificationToken = `
		INSERT INTO doctor_verification_tokens (token, expired_at, doctor_id)
		VALUES ($1, $2, $3);
	`

	qUpdateDoctorColl = `
		UPDATE doctors SET 
		name = $2,
		fee = $3,
		work_start_year = $4,
		doctor_specialists_id = $5,
	`

	qUpdateIsOnlineDoctor = `
		UPDATE doctors SET 
		is_online = $2,
		updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`
	qUpdateDoctorCommand = `
		updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
`

	qDeleteDoctor = `
		UPDATE doctors SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qCountTotalRows = `
		SELECT COUNT(*) OVER() AS total_rows  
	`

	qDoctorColl = `
		, d.id, d.name, d.email ,d.certificate, d.is_online, d.is_verified, d.is_google, d.fee, d.work_start_year, d.doctor_specialists_id, ds.name, d.profile_picture 
	`

	qDoctorCommands = `
		FROM doctors d
		LEFT JOIN doctor_specialists ds ON d.doctor_specialists_id = ds.id
		WHERE d.deleted_at IS NULL 
	`

	qUpadateDoctorPassword = `
	UPDATE doctors SET password = $2, updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL;
`

	qFindDoctorPasswordById = `
	SELECT id, password FROM doctors
	WHERE id = $1 AND deleted_at IS NULL
`

	// specialist
	qFindAllSpecialist = `
		SELECT id, name FROM doctor_specialists
		WHERE deleted_at IS NULL;
	`
)

const (
	qFindPharmacyManagerByEmail = `
		SELECT id, name, email, password, logo FROM pharmacy_managers
		WHERE email = $1 AND deleted_at IS NULL;
	`

	qPharmacyManagerColl = `
		, id, name, email, phone_number, logo 
	`
	qPharmacyManagerCommands = `
		FROM pharmacy_managers
		WHERE deleted_at IS NULL
	`

	qFindPharmacyMangerById = `
		SELECT id, name, email, phone_number, logo FROM pharmacy_managers
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qUpdatePharmacyManagerColl = `
		UPDATE pharmacy_managers SET
		name = $2,
		phone_number = $3,`

	qUpdatePharmacyManagerCommand = `
		updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qDeletePharmacyManagerById = `
		UPDATE pharmacy_managers SET
		deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`
	qCreateOnePharmacyManager = `
		INSERT INTO pharmacy_managers (name, email, password, phone_number, logo)
		VALUES ($1, $2, $3, $4, $5);
	`
)

const (
	qFindAdminByEmail = `
		SELECT id, name, email, password FROM admins
		WHERE email = $1 AND deleted_at IS NULL;
	`

	qCreateOneAdmin = `
		INSERT INTO admins (name, email, password) VALUES
		($1, $2, $3) RETURNING id;
	`

	qDeleteOneAdmin = `
		UPDATE admins SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL;
	`

	qFindOneAdminById = `
		SELECT id, name, email FROM admins
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qAdminColl = `
		, id, name, email
	`

	qAdminCommand = `
		FROM admins
		WHERE deleted_at IS NULL
	`
)

const (
	qFindPharmacyById = `
		SELECT 
		p.id, p.name, p.operational_hour, p.operational_day, p.pharmacist_name, p.pharmacist_license_number, p.pharmacist_phone_number,
		pm.id, pm.name, pm.email, pm.phone_number, pm.logo, pa.id, pa.pharmacy_id, pa.city, pa.province, pa.address, pa.district, pa.sub_district,
		pa.postal_code, ST_AsEWKT(pa.coordinate)
		FROM pharmacies p
		JOIN pharmacy_managers pm ON p.pharmacy_manager_id = pm.id
		JOIN pharmacy_addresses pa ON pa.pharmacy_id = p.id
		WHERE p.id = $1 AND p.deleted_at IS NULL;
	`

	qCreateOnePharmacy = `
		INSERT INTO pharmacies (name, pharmacy_manager_id, operational_hour, operational_day, pharmacist_name, pharmacist_license_number, pharmacist_phone_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	qUpdatePharmacy = `
		UPDATE pharmacies SET 
		name = $2,
		operational_hour = $3,
		operational_day = $4,
		pharmacist_name = $5,
		pharmacist_license_number = $6,
		pharmacist_phone_number = $7,
		updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qDeletePharmacyById = `
		UPDATE pharmacies SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL;
	`

	qPharmacyColl = `
		,p.id, p.name, p.operational_hour, p.operational_day, p.pharmacist_name, p.pharmacist_license_number, p.pharmacist_phone_number,
		pm.id, pm.name, pm.email, pm.phone_number, pm.logo, pa.id, pa.pharmacy_id, pa.city, pa.province, pa.address, pa.district, pa.sub_district,
		pa.postal_code, ST_AsEWKT(pa.coordinate)
	`

	qPharmacyCommands = `
		FROM pharmacies p
		JOIN pharmacy_managers pm ON p.pharmacy_manager_id = pm.id
		JOIN pharmacy_addresses pa ON pa.pharmacy_id = p.id
		WHERE p.pharmacy_manager_id = $1 AND p.deleted_at IS NULL
	`
)

const (
	qFindUserAddressByUserId = `
		SELECT id, user_id, city_id, city, province, address, district, sub_district, postal_code, ST_AsEWKT(coordinate), is_main
		FROM user_addresses
		WHERE user_id = $1 AND deleted_at IS NULL ORDER BY is_main DESC;
	`

	qCreateOneUserAddress = `
		INSERT INTO user_addresses(user_id, city, province, address, district, sub_district, postal_code, coordinate, is_main, city_id) VALUES
		($1, $2, $3, $4, $5, $6, $7, ST_MakePoint($8, $9), $10, $11)
		RETURNING id;
	`

	qUpdateUserAdrress = `
		UPDATE user_addresses SET
		city = $3,
		province = $4,
		address = $5,
		district = $6,
		sub_district = $7,
		postal_code = $8,
		coordinate = ST_MakePoint($9, $10),
		is_main = $11,
		city_id = $12,
		updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;
	`

	qDeleteUserAddress = `
		UPDATE user_addresses SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;
	`

	qFindUserAddressById = `
		SELECT id, user_id, city_id ,city, province, address, district, sub_district, postal_code, ST_AsEWKT(coordinate), is_main
		FROM user_addresses
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL ORDER BY is_main DESC;
	`

	qUpdateIsMain = `
		UPDATE user_addresses SET
		is_main = CASE WHEN id = $1 THEN NOT is_main
		ELSE is_main END WHERE id = $2;	
	`

	qUpdateIsMainFalse = `
		UPDATE user_addresses SET
		is_main = CASE WHEN id = $1 THEN true
		ELSE false END WHERE user_id = $2 AND deleted_at IS NULL;
	`

	qFindMainUserAddressByUserId = `
		SELECT id, user_id, city_id, city, province, address, district, sub_district, postal_code, ST_AsEWKT(coordinate), is_main
		FROM user_addresses
		WHERE user_id = $1 AND is_main = true AND deleted_at IS NULL;
	`
)

const (
	qFindGenderById = `
		SELECT id, name
		FROM genders
		WHERE id = $1 AND deleted_at IS NULL;
	`
)

const (
	qFindNearestPharmacy = `
		SELECT pharmacies.id
		FROM pharmacy_addresses
		JOIN pharmacies ON pharmacies.id = pharmacy_addresses.pharmacy_id
		JOIN pharmacy_products ON pharmacy_products.pharmacy_id = pharmacies.id
		WHERE ST_DWithin(pharmacy_addresses.coordinate, ST_MakePoint($1, $2)::geography, $3)
		AND pharmacies.deleted_at IS NULL
		GROUP BY pharmacy_addresses.coordinate, pharmacies.id
		ORDER BY SUM(pharmacy_products.total_stock) DESC, pharmacy_addresses.coordinate <-> ST_MakePoint($1, $2)::geography;
	`
	qFindPharmacyByProduct = `
		SELECT pharmacies.id, pharmacies.name, ST_DistanceSphere(ST_MakePoint(ST_X(coordinate::geometry), ST_Y(coordinate::geometry)),ST_MakePoint($1, $2))/1000 as distance
		FROM products
		JOIN pharmacy_products ON pharmacy_products.product_id = products.id 
		JOIN pharmacies ON pharmacies.id = pharmacy_products.pharmacy_id
		JOIN pharmacy_addresses ON pharmacy_addresses.pharmacy_id = pharmacies.id
		WHERE ST_DWithin(pharmacy_addresses.coordinate, ST_MakePoint($1, $2)::geography, $3) and products.id = $4
		AND pharmacy_products.deleted_at IS NULL
		GROUP BY pharmacy_addresses.coordinate, products.id, pharmacy_products.id, pharmacies.id, pharmacy_addresses.id, pharmacies.name
		ORDER BY pharmacy_addresses.coordinate <-> ST_MakePoint($1, $2)::geography;
	`

	qFindNearestPharmacyProductByProduct = `
		SELECT pharmacy_products.id
		FROM products
		JOIN pharmacy_products ON pharmacy_products.product_id = products.id 
		JOIN pharmacies ON pharmacies.id = pharmacy_products.pharmacy_id
		JOIN pharmacy_addresses ON pharmacy_addresses.pharmacy_id = pharmacies.id
		WHERE ST_DWithin(pharmacy_addresses.coordinate, ST_MakePoint($1, $2)::geography, $3) and products.id = $4
		AND pharmacy_products.deleted_at IS NULL
		GROUP BY pharmacy_addresses.coordinate, products.id, pharmacy_products.id, pharmacies.id, pharmacy_addresses.id, pharmacies.name
		ORDER BY pharmacy_addresses.coordinate <-> ST_MakePoint($1, $2)::geography
		LIMIT 1;
	`

	qNearestPharmacyProductColl = `
		, id, product_id, name, price, product_picture, selling_unit, slug_id
	`

	qNearestPharmacyProductCommands = `
		FROM (SELECT * from (
			SELECT pp.id, p.id AS product_id, p.name, pp.price, p.product_picture, p.selling_unit, p.slug_id FROM pharmacy_products pp 
			JOIN products p ON p.id = pp.product_id JOIN pharmacies ON pharmacies.id = pp.pharmacy_id
			JOIN pharmacy_addresses ON pharmacy_addresses.pharmacy_id = pharmacies.id
			JOIN product_categories pc ON pc.product_id = p.id
			WHERE ST_DWithin(pharmacy_addresses.coordinate, ST_MakePoint($1, $2)::geography, $3) 
	`

	qNearestPharmacyProductCommandsSecond = `
		AND pp.deleted_at IS NULL
			GROUP BY pharmacy_addresses.coordinate, p.id, pp.id, pharmacies.id, pharmacy_addresses.id, pharmacies.name
			ORDER BY SUM(pp.total_stock) DESC, pharmacy_addresses.coordinate <-> ST_MakePoint($1, $2)::geography) 
			ORDER BY row_number() OVER(PARTITION BY product_id) = 1 DESC
		FETCH FIRST 1 ROWS WITH TIES)
	`

	qFindNearestPharmacyMostBought = `
		SELECT pharmacies.id
		FROM pharmacy_addresses
		JOIN pharmacies ON pharmacies.id = pharmacy_addresses.pharmacy_id
		JOIN pharmacy_products ON pharmacy_products.pharmacy_id = pharmacies.id
		JOIN order_items ON order_items.pharmacy_product_id = pharmacy_products.id
		WHERE ST_DWithin(pharmacy_addresses.coordinate, ST_MakePoint($1, $2)::geography, $3)
		AND pharmacies.deleted_at IS NULL
		GROUP BY pharmacy_addresses.coordinate, pharmacies.id
		ORDER BY SUM(order_items.quantity) DESC, pharmacy_addresses.coordinate <-> ST_MakePoint($1, $2)::geography;
`
)

const (
	// pharmacy address
	qCreatePharmacyAddess = `
		INSERT INTO pharmacy_addresses(pharmacy_id, city_id, city, province, address, district, sub_district, postal_code, coordinate) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, ST_MakePoint($9, $10))
		RETURNING id;
	`

	qUpdatePharmacyAdrress = `
		UPDATE pharmacy_addresses SET
		city = $3,
		province = $4,
		address = $5,
		district = $6,
		sub_district = $7,
		postal_code = $8,
		coordinate = ST_MakePoint($9, $10),
		updated_at = NOW()
		WHERE id = $1 AND pharmacy_id = $2 AND deleted_at IS NULL;
	`

	qDeletePharmacyAddress = `
		UPDATE pharmacy_addresses SET deleted_at = NOW()
		WHERE pharmacy_id = $1 AND deleted_at IS NULL;
	`

	// pharmacy product
	qFindPharmacyProduct = `
		SELECT products.product_picture, products.name, pharmacy_products.price, products.selling_unit, 
		products.slug_id, pharmacy_products.id as pharmacy_product_id, products.id as product_id, COUNT(*) OVER () as total
		FROM pharmacy_products 
		JOIN products ON products.id = pharmacy_products.product_id
		WHERE pharmacy_products.pharmacy_id = $1
		AND pharmacy_products.deleted_at IS NULL
		OFFSET $2
		LIMIT $3;
	`
	qFindPharmacyProductByCategory = `
		SELECT products.product_picture, products.name, pharmacy_products.price, products.selling_unit, products.slug_id, pharmacy_products.id as pharmacy_product_id, products.id as product_id, COUNT(*) OVER () as total
		FROM pharmacy_products 
		JOIN products ON products.id = pharmacy_products.product_id
		JOIN product_categories ON product_categories.product_id = products.id
		WHERE pharmacy_products.pharmacy_id = $1 AND product_categories.category_id = $4 AND pharmacy_products.deleted_at IS NULL AND pharmacy_products.is_available = true
		OFFSET $2
		LIMIT $3;
	`

	qFindPharmacyProductByPharmacyId = `
	SELECT products.product_picture, products.name, pharmacy_products.price, products.selling_unit, 
	products.slug_id, pharmacy_products.id as pharmacy_product_id, products.id as product_id
	FROM pharmacy_products 
	JOIN products ON products.id = pharmacy_products.product_id
	WHERE pharmacy_products.pharmacy_id = $1
	AND pharmacy_products.deleted_at IS NULL
`
)

const (
	qCreatePharmacyProduct = `
		INSERT INTO pharmacy_products (price, total_stock, is_available, product_id, pharmacy_id, deleted_at) VALUES
		($1, $2, $3, $4, $5, NULL) RETURNING id;
	`

	qUpdatePharmacyProduct = `
		UPDATE pharmacy_products SET
		price = $2, total_stock = $3, is_available = $4, product_id = $5, pharmacy_id = $6, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qDeletePharmacyProduct = `
		UPDATE pharmacy_products SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qFindOnePharmacyProduct = `
		SELECT id, pharmacy_id, product_id, total_stock FROM pharmacy_products
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qFindOneByPharmacyAndProductId = `
	SELECT pp.id, pp.total_stock, pp.is_available,
	ph.id, ph.name, 
	pd.id, pd.name, pm.id,
	pp.deleted_at 
	FROM pharmacy_products pp
	JOIN pharmacies ph ON pp.pharmacy_id = ph.id
	JOIN products pd ON pp.product_id = pd.id 
	JOIN pharmacy_managers pm ON ph.pharmacy_manager_id = pm.id
	WHERE pharmacy_id = $1 AND product_id = $2 AND pp.deleted_at IS NULL;
`
	qFindOnePharmacyProductById = `
		SELECT 
		pp.id, pp.total_stock, pp.is_available, pp.price,
		pd.id, pd.name, pd.content, pd.description, pd.unit_in_pack, pd.selling_unit,
		pd.weight, pd.height, pd.length, pd.width, pd.product_picture, pd.slug_id, pc.name,
		pf.name, m.name,
		pp.deleted_at 
		FROM pharmacy_products pp
		JOIN products pd ON pp.product_id = pd.id 
		JOIN product_forms pf ON pd.product_form_id = pf.id
		JOIN product_classifications pc ON pd.product_classification_id = pc.id
		JOIN manufactures m ON pd.manufacture_id = m.id
		WHERE pp.id = $1 AND pp.deleted_at IS NULL;
`

	qUpdatePharmacyProductStock = `
		UPDATE pharmacy_products SET 
		total_stock = $3, updated_at = NOW()
		WHERE pharmacy_id = $1 AND product_id = $2 AND deleted_at IS NULL; 
	`
)

const (
	qDeletedAllPharmacyProduct = `
		UPDATE pharmacy_products SET
		deleted_at = NOW() WHERE pharmacy_id = $1 AND deleted_at IS NULL;
	`

	qPharmacyProductByPharmacyIdColl = `
		, pp.id, pp.total_stock, pp.is_available, pp.price,
		pd.id, pd.name, pd.content, pd.description, pd.unit_in_pack, pd.selling_unit,
		pd.weight, pd.height, pd.length, pd.width, pd.product_picture, pd.slug_id, pc.name,
		pf.name, m.name
`

	qPharmacyProductByPharmacyIdCommand = `
	FROM pharmacy_products pp
	JOIN products pd ON pp.product_id = pd.id 
	JOIN product_forms pf ON pd.product_form_id = pf.id
	JOIN product_classifications pc ON pd.product_classification_id = pc.id
	JOIN manufactures m ON pd.manufacture_id = m.id
	JOIN pharmacies p ON pp.pharmacy_id = p.id
	WHERE pp.pharmacy_id = $1 AND p.pharmacy_manager_id = $2 AND pp.deleted_at IS NULL
	`

	// category
	qCategoryColl = `
		, id, name
	`

	qCategoryCommands = `
		FROM categories
		WHERE deleted_at IS NULL
	`

	qFindOneCategoryById = `
		SELECT id, name FROM categories
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qCreateOneCategory = `
		INSERT INTO categories(name) VALUES
		($1) RETURNING id;
	`

	qUpdateOneCategory = `
		UPDATE categories SET
		name = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qDeleteOneCategory = `
		UPDATE categories SET
		deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`
)

const (
	qFindPharmacyProductDetail = `
		SELECT pp.id, p.name, p.generic_name, p.content, p.description, p.unit_in_pack, p.selling_unit, p.weight, p.height, p.length, p.width, 
		p.product_picture, p.slug_id, pf.name AS product_form, pc.name AS product_classification, m.name AS manufacture, pp.price, pp.total_stock
		FROM pharmacy_products pp
		JOIN products p ON p.id = pp.product_id 
		JOIN product_forms pf ON pf.id = p.product_form_id 
		JOIN product_classifications pc ON pc.id = p.product_classification_id 
		JOIN manufactures m ON m.id = p.manufacture_id 
		WHERE pp.id = $1 AND pp.deleted_at IS NULL;
	`
)

const (
	qFindProductCategories = `
		SELECT c.id, c.name
		FROM products p 
		JOIN product_categories pc ON pc.product_id = p.id
		JOIN categories c ON c.id = pc.category_id 
		WHERE p.id = $1;
	`
)

const (
	qFindPharmacyOfficialShippingMethod = `
		SELECT osm.id, osm.name, osm.fee
		FROM pharmacies_shipping_methods psm 
		JOIN official_shipping_methods osm ON osm.id = psm.official_id 
		WHERE psm.pharmacy_id = $1 AND psm.deleted_at IS NULL;
	`

	qFindPharmacyNonOfficialShippingMethod = `
		SELECT nosm.id, nosm.name, nosm.courier, nosm.service, nosm.description 
		FROM pharmacies_shipping_methods psm 
		JOIN non_official_shipping_methods nosm ON nosm.id = psm.non_official_id 
		WHERE psm.pharmacy_id = $1 AND psm.deleted_at IS NULL;
	`

	qCreatePharmacyOfficialShippingMethod = `
		INSERT INTO pharmacies_shipping_methods(official_id, pharmacy_id) VALUES
		($1, $2) RETURNING id
	`

	qCreatePharmacyNonOfficialShippingMethod = `
		INSERT INTO pharmacies_shipping_methods(non_official_id, pharmacy_id) VALUES
		($1, $2) RETURNING id
	`

	qUpdatedOfficialPharmacyShippingMethod = `
		UPDATE pharmaies_shipping_methods SET
		official_id = $2, pharmacy_id = $3, updated_at = NOW()
		WHERE id = $1 and deleted_at IS NULL
	`

	qUpdatedNonOfficialPharmacyShippingMethod = `
		UPDATE pharmaies_shipping_methods SET
		non_official_id = $2, pharmacy_id = $3, updated_at = NOW()
		WHERE id = $1 and deleted_at IS NULL
	`

	qDeletedAllPharmacyShippingMethod = `
		UPDATE pharmacies_shipping_methods SET
		deleted_at = NOW() WHERE pharmacy_id = $1 and deleted_at IS NULL;
	`

	qDeletedShippingMethodByPharmacyId = `
		DELETE FROM pharmacies_shipping_methods
		WHERE pharmacy_id = $1 AND deleted_at IS NULL
	`
)

const (
	qCreateOneProduct = `
		INSERT INTO products(
		name, generic_name, content, description, unit_in_pack, selling_unit, weight, height, length,
		width, product_picture, slug_id, product_form_id, product_classification_id, manufacture_id
		) 
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id;
	`

	qUpdateOneProductColl = `
		UPDATE products SET
		name = $2, 
		generic_name = $3, 
		content = $4, 
		description = $5, 
		unit_in_pack = $6, 
		selling_unit = $7, 
		weight = $8, 
		height = $9, 
		length = $10, 
		width = $11, 
		slug_id = $12, 
		product_form_id = $13, 
		product_classification_id = $14,
		manufacture_id = 	$15,
	`

	qUpdateOneProductCommand = `
		updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	qDeleteOneProduct = `
		UPDATE products SET
		deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qFindOneProductById = `
		SELECT 
		p.id, p.name, p.generic_name, p.content, p.description, p.unit_in_pack, p.selling_unit, p.weight, p.height, p.length,
		p.width, p.product_picture, p.slug_id, pf.id, pf.name AS product_form, pc.id, pc.name AS product_classification, m.id, m.name AS manufacture
		FROM products p
		JOIN product_forms pf ON pf.id = p.product_form_id
		JOIN product_classifications pc ON pc.id = p.product_classification_id
		JOIN manufactures m ON m.id = p.manufacture_id
    WHERE p.id = $1 AND p.deleted_at IS NULL
	`

	qProductColl = `
		, p.id, p.name, p.generic_name, p.content, p.description, p.unit_in_pack, p.selling_unit, p.weight, p.height, p.length,
		p.width, p.product_picture, p.slug_id, pf.id, pf.name AS product_form, pc.id, pc.name AS product_classification, m.id, m.name AS manufacture
	`

	qProductCommands = `
		FROM products AS p
		JOIN product_forms AS pf ON pf.id = p.product_form_id
		JOIN product_classifications AS pc ON pc.id = p.product_classification_id
		JOIN manufactures AS m ON m.id = p.manufacture_id
		WHERE p.deleted_at IS NULL
	`
)

const (
	qCreateOneProductCategory = `
		INSERT INTO product_categories(product_id, category_id) VALUES
		($1, $2) RETURNING id;
	`

	qUpdateOneProductCategoryByProductId = `
		UPDATE product_categories SET
		category_id = $2,
		updated_at = NOW()
		WHERE product_id = $1 AND deleted_at IS NULL;
	`

	qFindAllProductCategoryByProductId = `
		SELECT id, product_id, category_id FROM product_categories
		WHERE product_id = $1 AND deleted_at IS NULL;
	`

	qDeleteProductCategoryByProductId = `
		UPDATE product_categories SET
		deleted_at = NOW()
		WHERE product_id = $1 AND deleted_at IS NULL;
	`
)

const (
	qFindAllProductForm = `
		SELECT id, name FROM product_forms WHERE deleted_at IS NULL
	`

	qFindAllClassification = `
		SELECT id, name FROM product_classifications WHERE deleted_at IS NULL
	`

	qFindAllManufacture = `
		SELECT id, name FROM manufactures WHERE deleted_at IS NULL
	`
)

const (
	qCreateCart = `
		INSERT INTO cart_items (quantity, user_id, pharmacy_product_id) VALUES
		($1, $2, $3);
	`
	qIncreaseCartQuantity = `
		UPDATE cart_items SET
		quantity = quantity + $1,
		updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL;
	`
	qDecreaseCartQuantity = `
		UPDATE cart_items SET
		quantity = quantity - $1,
		updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING id, quantity;
	`
	qDeleteCartItem = `
		UPDATE cart_items SET
		deleted_at = NOW(),
		quantity = 0
		WHERE id = $1 AND deleted_at IS NULL;
	`
	qGetUserCart = `
		SELECT quantity, user_id, pharmacy_product_id
		FROM cart_items
		WHERE user_id = $1 AND deleted_at IS NULL;
	`
	qFindOneCartItem = `
		SELECT id, quantity
		FROM cart_items
		WHERE user_id = $1 AND pharmacy_product_id = $2 AND deleted_at IS NULL;
	`
	qFindCartItem = `
		SELECT * from (select ci.id AS cart_item_id, ci.updated_at, p.name AS product_name, p.product_picture, p.selling_unit, pp.price, ci.quantity, p2.name AS pharmacy_name, p2.id AS pharmacy_id, p.slug_id, pp.total_stock, pp.is_available, p.weight
		FROM cart_items ci 
		JOIN pharmacy_products pp ON pp.id = ci.pharmacy_product_id 
		JOIN products p ON p.id = pp.product_id 
		JOIN pharmacies p2 ON p2.id = pp.pharmacy_id
		WHERE ci.user_id = $1 AND ci.deleted_at IS NULL) AS tab
		ORDER BY MAX(tab.updated_at) OVER (partition BY tab.pharmacy_id) DESC, tab.updated_at DESC;
	`
	qFindPharmacyIdByCartId = `
		SELECT pp.pharmacy_id, ci.quantity, pp.id
		FROM cart_items ci 
		JOIN pharmacy_products pp ON pp.id = ci.pharmacy_product_id 
		WHERE ci.id = $1 AND ci.deleted_at IS NULL;
	`

	// location
	qFindProvinces       = `SELECT id, name FROM provinces;`
	qFindCities          = `SELECT id, name, type FROM cities WHERE province_id = $1`
	qFindDistricts       = `SELECT id, name FROM districts WHERE city_id = $1`
	qFindSubDistricts    = `SELECT id, name, postal_code, ST_AsEWKT(coordinate) FROM sub_districts WHERE district_id = $1`
	qFindProvinceAndCity = `
		SELECT p.id AS province_id, c.id AS city_id FROM provinces p
		JOIN cities c ON p.id = c.province_id
		WHERE c.name = $1 AND c.type = $2;
	`
	dFindDistrictAndSubDistrict = `
		SELECT district_id, sub_district_id, postal_code
		FROM (
    		SELECT d.id AS district_id, sd.id AS sub_district_id, sd.postal_code, 1 AS order_column
    		FROM cities c
    		JOIN districts d ON c.id = d.city_id
    		JOIN sub_districts sd ON d.id = sd.district_id
    		WHERE c.id = $1 AND sd.name = CONCAT(--sub-district-concat--)
    		UNION ALL
    		SELECT d.id AS district_id, sd.id AS sub_district_id, sd.postal_code, 2 AS order_column
    		FROM cities c
    		JOIN districts d ON c.id = d.city_id
    		JOIN sub_districts sd ON d.id = sd.district_id
    		WHERE c.id = $1 AND (--sub-district-ilike--)
		) combined_results
		ORDER BY order_column;
	`

	// stock history
	qCreateOneStockHistory = `
		INSERT INTO stock_histories(pharmacy_product_id, pharmacy_id, quantity, description) VALUES
		($1, $2, $3, $4) RETURNING id
	`

	qUpdateStockHistory = `
		UPDATE stock_histories SET
		pharmacy_product_id = $3, pharmacy_id = $4, quantity = $5, description = $6,
		updated_at = NOW() WHERE pharmacy_product_id = $1 AND pharmacy_id = $2 AND deleted_at IS NULL;
	`

	qFindStockHistoryByPharmacyIdAndPharmacyProductId = `
		SELECT sh.id, sh.quantity, sh.description,
		pp.id, pp.price,
		pd.id, pd.name,
		p.id, p.name
		FROM stock_histories
		JOIN pharmacy_products pp ON sh.pharmacy_product_id = pp.id
		JOIN products pd ON pp.product_id = pd.id 
		JOIN pharmacies p ON pp.pharmacy_id = p.id 
		WHERE sh.pharmacy_id = $1 AND sh.pharmacy_product_id = $2 sh.deleted_at IS NULL;
	`

	qDeleteAllStockHistoryByPharmacyId = `
		UPDATE stock_histories SET
		deleted_at = NOW() WHERE pharmacy_id = $1 AND deleted_at IS NULL;
	`

	qFindStockHistoriesByPharmacyIdColl = `
			, sh.id, sh.quantity, sh.description, 
		pp.id, pp.price,
		pd.id, pd.name, 
		p.id, p.name, p.pharmacy_manager_id
	`

	qFindStockHistoriesByPharmacyIdCommand = `
		FROM stock_histories sh
		JOIN pharmacy_products pp ON sh.pharmacy_product_id = pp.id
		JOIN products pd ON pp.product_id = pd.id 
		JOIN pharmacies p ON sh.pharmacy_id = p.id
		WHERE sh.pharmacy_id = $1 AND p.pharmacy_manager_id = $2 AND sh.deleted_at IS NULL
	`

	//shipping cost
	qFindUserPharmacyDistance = `
		SELECT ST_DistanceSphere(ST_MakePoint(ST_X(ua.coordinate::geometry), ST_Y(ua.coordinate::geometry)),ST_MakePoint(ST_X(pa.coordinate::geometry), ST_Y(pa.coordinate::geometry)))/1000 AS distance
		FROM user_addresses ua, pharmacy_addresses pa
		WHERE ua.id = $1 AND pa.pharmacy_id = $2 AND ua.deleted_at IS NULL AND pa.deleted_at IS NULL;
	`
	qFindOfficialShippingFee = `
		SELECT fee
		FROM official_shipping_methods osm 
		WHERE osm.id = $1 AND deleted_at IS NULL;
	`
	qFindNonOfficialShippingMethod = `
		SELECT courier, service 
		FROM non_official_shipping_methods nosm 
		WHERE nosm.id = $1 AND deleted_at IS NULL;
	`

	//pharmacy address
	qFindPharmacyAddress = `
		SELECT city_id
		FROM pharmacy_addresses
		WHERE pharmacy_addresses.pharmacy_id = $1 AND deleted_at IS NULL;
	`

	// order
	qCreateOrder = `
		INSERT INTO orders (order_number, total_price, payment_deadline, shipping_fee, shipping_method, user_address_id, order_status_id, pharmacy_id)
		VALUES ($1, $2, $3, $4, $5, $6, (select id from order_statuses where name = $7), $8)
		RETURNING id, payment_deadline
	`
	qCreateOrderItem = `
		INSERT INTO order_items (quantity, order_id, pharmacy_product_id)
		VALUES %s
	`
	qCartItemsBulkDelete = `
		UPDATE cart_items  
		SET deleted_at = NOW()
		FROM
			( VALUES
				%s
			) AS nv (id)
		WHERE cart_items.id = CAST (nv.id AS BIGINT);
	`
	qFindUserOrder = `
		SELECT o.id, o.order_number, o.total_price, o.payment_proof, o.payment_deadline, o.shipping_fee, o.shipping_method, os.name AS status, p.name AS pharmacy_name, ua.city, ua.province, ua.address, ua.district, ua.sub_district, ua.postal_code, COUNT(*) OVER () as total
		FROM orders o 
		JOIN user_addresses ua ON ua.id = o.user_address_id 
		JOIN order_statuses os ON os.id = o.order_status_id 
		JOIN pharmacies p ON p.id = o.pharmacy_id 
		WHERE ua.user_id = $1 AND LOWER(os.name) = lower($2) AND o.deleted_at IS NULL 
		ORDER BY o.created_at DESC
		OFFSET $3
		LIMIT $4;
	`
	qFindOrderDetail = `
		SELECT o.id, o.order_number, o.total_price, o.payment_proof, o.payment_deadline, o.shipping_fee, o.shipping_method, os.name AS status, p.name AS pharmacy_name, ua.city, ua.province, ua.address, ua.district, ua.sub_district, ua.postal_code
		FROM orders o 
		JOIN user_addresses ua ON ua.id = o.user_address_id 
		JOIN order_statuses os ON os.id = o.order_status_id 
		JOIN pharmacies p ON p.id = o.pharmacy_id 
		WHERE o.id = $1 AND o.deleted_at IS NULL;
	`
	qFindOrderItems = `
		SELECT p.name, p.selling_unit, pp.price, oi.quantity, p.product_picture, pp.id
		FROM order_items oi 
		JOIN pharmacy_products pp ON pp.id = oi.pharmacy_product_id 
		JOIN products p ON p.id = pp.product_id 
		WHERE oi.order_id = $1 AND oi.deleted_at IS NULL;
	`
	qFindUserOrderByPharmacyManager = `
		SELECT o.id, o.order_number, o.total_price, o.payment_proof, o.payment_deadline, o.shipping_fee, o.shipping_method, os.name AS status, p.name AS pharmacy_name, ua.city, ua.province, ua.address, ua.district, ua.sub_district, ua.postal_code, pa.city, pa.province, pa.address, pa.district, pa.sub_district, pa.postal_code, COUNT(*) OVER () as total, u.name AS user_name, u.email AS user_email
		FROM orders o 
		JOIN user_addresses ua ON ua.id = o.user_address_id 
		JOIN order_statuses os ON os.id = o.order_status_id 
		JOIN pharmacies p ON p.id = o.pharmacy_id
		JOIN pharmacy_addresses pa ON pa.pharmacy_id = p.id
		JOIN users u ON u.id = ua.user_id
		WHERE p.pharmacy_manager_id = $1 AND LOWER(os.name) = lower($2) AND o.deleted_at IS NULL %s
		ORDER BY o.created_at DESC
		OFFSET $3
		LIMIT $4;
	`
	qFindUserOrderByAdmin = `
		SELECT o.id, o.order_number, o.total_price, o.payment_proof, o.payment_deadline, o.shipping_fee, o.shipping_method, os.name AS status, p.name AS pharmacy_name, ua.city, ua.province, ua.address, ua.district, ua.sub_district, ua.postal_code, pa.city, pa.province, pa.address, pa.district, pa.sub_district, pa.postal_code, COUNT(*) OVER () as total, u.name AS user_name, u.email AS user_email 
		FROM orders o 
		JOIN user_addresses ua ON ua.id = o.user_address_id 
		JOIN order_statuses os ON os.id = o.order_status_id 
		JOIN pharmacies p ON p.id = o.pharmacy_id
		JOIN pharmacy_addresses pa ON pa.pharmacy_id = p.id
		JOIN users u ON u.id = ua.user_id
		WHERE LOWER(os.name) = lower($1) AND o.deleted_at IS NULL %s
		ORDER BY o.created_at DESC
		OFFSET $2
		LIMIT $3;
	`
	qUpdateOrderStatus = `
		UPDATE orders SET
		order_status_id = (SELECT id FROM order_statuses WHERE LOWER(name) = lower($1) AND deleted_at IS NULL),
		updated_at = NOW()
		FROM pharmacies, user_addresses
		WHERE orders.id = $2 AND orders.deleted_at IS NULL %s
	`
	qFindOrderStatus = `
		SELECT os.name, o.payment_proof, o.pharmacy_id
		FROM orders o 
		JOIN order_statuses os ON os.id = o.order_status_id 
		WHERE o.id = $1 AND o.deleted_at IS NULL;
	`
	qUploadPaymentProof = `
		UPDATE orders SET 
		payment_proof = $1,
		updated_at = NOW()
		FROM user_addresses
		WHERE orders.id = $2 AND orders.deleted_at IS NULL AND user_addresses.user_id = orders.user_address_id AND user_addresses.user_id = $3;
	`

	// stock
	qIncreaseStock = `
		UPDATE pharmacy_products SET 
		total_stock = total_stock + $1
		WHERE id = $2;
	`
	qDecreaseStock = `
		UPDATE pharmacy_products SET 
		total_stock = total_stock - $1
		WHERE id = $2;
	`
	qLockPharmacyProductRow = `
		SELECT * FROM pharmacy_products WHERE id = $1 FOR UPDATE
	`
)

const (
	qFindConsultationById = `
		SELECT c.id, d.id, d.name, d.profile_picture, d.is_online, s.id, s.name, u.id, u.name, u.profile_picture, g.id, g.name, c.patient_name, c.patient_birth_date, c.certificate_url, c.prescription_url, c.ended_at, c.created_at FROM consultations c
		JOIN doctors d ON c.doctor_id = d.id 
		JOIN doctor_specialists s ON d.doctor_specialists_id = s.id
		JOIN users u ON c.user_id = u.id
		JOIN genders g ON c.patient_gender_id = g.id
		WHERE c.id = $1 AND c.deleted_at IS NULL;
	`

	qConsultationColl = `
		, c.id, d.id, d.name, d.profile_picture, d.is_online, s.id, s.name, u.id, u.name, u.profile_picture, g.id, g.name, c.patient_name, c.patient_birth_date, c.certificate_url, c.prescription_url, c.ended_at, c.created_at
	`

	qConsultationCommands = `
		FROM consultations c
		JOIN doctors d ON c.doctor_id = d.id 
		JOIN doctor_specialists s ON d.doctor_specialists_id = s.id
		JOIN users u ON c.user_id = u.id
		JOIN genders g ON c.patient_gender_id = g.id
	`

	qConsultationByUserIdCommands = `
		WHERE c.user_id = $1 AND c.deleted_at IS NULL
	`
	qConsultationByDoctorIdCommands = `
		WHERE c.doctor_id = $1 AND c.deleted_at IS NULL
	`

	qCreateOneConsultation = `
		INSERT INTO consultations(doctor_id, user_id, patient_gender_id, patient_name, patient_birth_date) VALUES
		($1, $2, $3, $4, $5) RETURNING id;
	`

	qUpdateEndedAtConsultation = `
		UPDATE consultations SET
		ended_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qUpdateCertificateConsultation = `
		UPDATE consultations SET
		certificate_url = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qCreatePrescriptionItems = `
		INSERT INTO prescription_items (consultation_id, product_id, quantity) VALUES  
	`

	qUpdatePrescriptionConsultation = `
		UPDATE consultations SET
		prescription_url = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	qPrescriptionProducts = `
		SELECT product_id, quantity FROM prescription_items 
		WHERE consultation_id = $1 AND deleted_at IS NULL;
	`

	qFindAllChatByConsultationId = `
		SELECT id, is_from_user, content, type, created_at FROM chats 
		WHERE consultation_id = $1 AND deleted_at IS NULL;
	`

	qCreateOneChat = `
		INSERT INTO chats (consultation_id, is_from_user, content, type) 
		VALUES ($1, $2, $3, $4) RETURNING id;
	`
	// user reset password
	qCreateUserResetPassword = `
		INSERT INTO user_reset_password_tokens(token, user_id, expired_at) VALUES
		($1, $2, $3)
		RETURNING id;
	`

	qFindUserResetPasswordByToken = `
		SELECT id, token, user_id, expired_at FROM user_reset_password_tokens
		WHERE token = $1 AND deleted_at IS NULL;
	`

	qFindUserResetPasswordByUserId = `
		SELECT id, token, user_id, expired_at FROM user_reset_password_tokens
		WHERE user_id = $1 AND deleted_at IS NULL;
	`

	qDeleteUserResetPassword = `
		UPDATE user_reset_password_tokens SET
		deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL;
	`

	// pharmacy reset password
	qCreateDoctorResetPassword = `
		INSERT INTO doctor_reset_password_tokens(token, doctor_id, expired_at) VALUES
		($1, $2, $3)
		RETURNING id;
	`

	qFindDoctorResetPasswordByToken = `
		SELECT id, token, doctor_id, expired_at FROM doctor_reset_password_tokens
		WHERE token = $1 AND deleted_at IS NULL;
	`

	qFindDoctorResetPasswordByDoctorId = `
		SELECT id, token, doctor_id, expired_at FROM doctor_reset_password_tokens
		WHERE doctor_id = $1 AND deleted_at IS NULL;
	`

	qDeleteDoctorResetPassword = `
		UPDATE doctor_reset_password_tokens SET
		deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL;
	`
)

const (
	qFindAllMutationStatus = `
		SELECT id, name FROM mutation_statuses
		WHERE deleted_at IS NULL;
	`
)

const (
	qCreateOneStockTransfer = `
		INSERT INTO stock_transfer_requests 
		(pharmacy_sender_id, pharmacy_receiver_id, mutation_status_id, product_id, quantity)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, pharmacy_sender_id, pharmacy_receiver_id, mutation_status_id, product_id, quantity;
	`

	qStockTrasnferColl = `
	, st.id, phs.id, phs.name, phr.id, phr.name, 
	ms.id, ms.name, pd.id, pd.name, st.quantity, st.updated_at
		`

	qStockTrasnferCommands = `
		FROM stock_transfer_requests st
		JOIN pharmacies phs ON st.pharmacy_sender_id = phs.id
		JOIN pharmacies phr ON st.pharmacy_receiver_id = phr.id
		JOIN mutation_statuses ms ON st.mutation_status_id = ms.id
		JOIN products pd ON st.product_id = pd.id
		WHERE st.deleted_at IS NULL
	`

	qUpdateOneMutationStatusId = `
		UPDATE stock_transfer_requests SET mutation_status_id = $2 WHERE id = $1 AND deleted_at IS NULL;
	`

	qFindOneStockTransfer = `
		SELECT st.id, st.pharmacy_sender_id, st.pharmacy_receiver_id, st.mutation_status_id, st.product_id, st.quantity
		FROM stock_transfer_requests st WHERE st.id = $1 AND st.deleted_at IS NULL
	`
)

const (
	qStockHistoryReportsColl = `
		,COALESCE(SUM(CASE WHEN sh.quantity > 0 THEN sh.quantity ELSE 0 END), 0) AS total_addition,
		COALESCE(SUM(CASE WHEN sh.quantity < 0 THEN sh.quantity ELSE 0 END), 0) AS total_deduction,
		pp.total_stock AS final_stock,
		pd.id,
		pd.name,
		ph.id,
		ph.name,
		DATE_TRUNC('month', sh.created_at) AS month,
		DATE_TRUNC('year', sh.created_at) AS year
	`

	qStockHistoryReportsCommand = `
		FROM stock_histories sh
		JOIN pharmacy_products pp ON sh.pharmacy_product_id = pp.id
		JOIN pharmacies ph ON pp.pharmacy_id = ph.id
		JOIN products pd ON pp.product_id = pd.id
		WHERE
		DATE_TRUNC('month', sh.created_at) >= DATE_TRUNC('month', CURRENT_DATE)
		AND DATE_TRUNC('month', sh.created_at) <= DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '12 month'
		AND sh.deleted_at IS NULL
	`

	qStockHistoryReportsGroup = `
		GROUP BY
		sh.pharmacy_product_id, pp.total_stock, month, pd.id, ph.id, month, year
	`
)

const (
	qSalesReportColl = `
		, ph.id, ph.name,
		pd.id, pd.name, pd.content, pd.description, pd.unit_in_pack, pd.selling_unit,
		pd.weight, pd.height, pd.length, pd.width, pd.product_picture, pd.slug_id,
		pc.name, pf.name, m.name,
		SUM(oi.quantity * pp.price) AS total_sales_amount,
		SUM(oi.quantity) AS total_quantity_sold,
		DATE_TRUNC('month', oi.created_at) AS month,
		DATE_TRUNC('year', oi.created_at) AS year
	`
	qSalesReportCommand = `
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		JOIN pharmacy_products pp ON oi.pharmacy_product_id = pp.id
		JOIN pharmacies ph ON pp.pharmacy_id = ph.id
		JOIN products pd ON pp.product_id = pd.id
		JOIN product_forms pf ON pd.product_form_id = pf.id
		JOIN product_classifications pc ON pd.product_classification_id = pc.id
		JOIN manufactures m ON pd.manufacture_id = m.id
		WHERE oi.deleted_at IS NULL AND
    DATE_TRUNC('month', oi.created_at) >= DATE_TRUNC('month', CURRENT_DATE) AND
    DATE_TRUNC('month', oi.created_at) <= DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '12 month'
	`

	qSalesReportGroup = `
		GROUP BY
		ph.id, ph.name, pd.id, pd.name, month, year, oi.id, pc.id, pf.id, m.id
	`
)

const (
	qSalesReportCategoryColl = `
		, c.id,
		c.name,
		EXTRACT(MONTH FROM o.created_at) AS month,
		EXTRACT(YEAR FROM o.created_at) AS year,
		SUM(oi.quantity)
	`

	qSalesReportCategoryCommand = `
		FROM orders o
		JOIN order_items oi ON o.id = oi.order_id
		JOIN pharmacy_products pp ON oi.pharmacy_product_id = pp.id
		JOIN products p ON pp.product_id = p.id
		JOIN product_categories pc ON p.id = pc.product_id
		JOIN categories c ON pc.category_id = c.id
		WHERE o.deleted_at IS NULL
	`

	qSalesReportCategoryGroup = `
		GROUP BY oi.id, c.id, c.name, EXTRACT(MONTH FROM o.created_at), EXTRACT(YEAR FROM o.created_at)
	`
)

const (
	qFindMostBought = `
	SELECT
    pp.id AS pharmacy_product_id,
    pp.price as price,
    pd.selling_unit as selling_unit,
    pd.slug_id as slug,
    pd.id AS product_id,
    pd.name AS product_name,
    pd.product_picture AS product_picture,
		DATE_TRUNC('day', oi.created_at) AS day,
    SUM(oi.quantity) as total_quantity_sold,
		COUNT(*) OVER ()
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		JOIN pharmacy_products pp ON oi.pharmacy_product_id = pp.id
		JOIN pharmacies ph ON pp.pharmacy_id = ph.id
		JOIN products pd ON pp.product_id = pd.id
		WHERE ph.id = $1
		AND DATE_TRUNC('day', oi.created_at) >= DATE_TRUNC('day', CURRENT_DATE) AND
		DATE_TRUNC('day', oi.created_at) <= DATE_TRUNC('day', CURRENT_DATE) + INTERVAL '2 day'
		GROUP BY
		pp.id, oi.created_at, ph.id, ph.name, pd.id, pd.name, EXTRACT(day FROM oi.created_at)
		ORDER BY day
				OFFSET $2
				LIMIT $3
			`
)
