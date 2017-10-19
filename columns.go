package store

const (
	BarColumns = "id, cafeteria_id, name, create_time, update_time, creator, operator, removed, version"

	CafeteriaColumns = "id, name, create_time, update_time, creator, operator, removed, version"

	CafeteriaUserColumns = "id, cafeteria_id, user_id, title, authorities, create_time, update_time, creator, operator, removed, version"

	CashRegisterColumns = "id, cafeteria_id, bar_id, name, create_time, update_time, creator, operator, removed, version"

	ContainerColumns = "id, cafeteria_id, bar_id, name, create_time, update_time, creator, operator, removed, version"

	InteractiveWindowColumns = "id, cafeteria_id, name, create_time, update_time, creator, operator, removed, version"

	LoginCredentialColumns = "id, user_id, login_type, credential, create_time, update_time, creator, operator, removed, version"

	ShopColumns = "id, cafeteria_id, window_id, name, restaurant_id, create_time, update_time, creator, operator, removed, version"

	UserInfoColumns = "id, username, password, salt, real_name, create_time, update_time, creator, operator, removed, version"
)
