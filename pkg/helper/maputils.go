package helper

type Keyed interface {
	GetKey() string
}

// из slice в map[]
func ToMapByKey[T Keyed](list []T) map[string]T {
	result := make(map[string]T, len(list))
	for _, item := range list {
		result[item.GetKey()] = item
	}
	return result
}

type ConvertibleToEntity[Entity any] interface {
	ToEntity() Entity
	GetKey() string
}

// FilterNewEntities фильтрует DTO, оставляя только те, ключи которых нет среди существующих сущностей, и конвертирует их в сущности.
func FilterNewEntities[DTO ConvertibleToEntity[Entity], Entity Keyed](dtos []DTO, ents []Entity) []Entity {
	existing := make(map[string]struct{}, len(ents))
	for _, ent := range ents {
		existing[ent.GetKey()] = struct{}{}
	}

	result := make([]Entity, 0, len(dtos))
	for _, dto := range dtos {
		if _, ok := existing[dto.GetKey()]; !ok {
			result = append(result, dto.ToEntity())
		}
	}
	return result
}

// CollectKeys принимает срез сущностей T (которые реализуют Keyed) и возвращает срез ключей из всех сущностей.
func CollectKeys[T Keyed](ents []T) []string {
	keys := make([]string, 0, len(ents))
	for _, ent := range ents {
		keys = append(keys, ent.GetKey())
	}
	return keys
}

// получаем ключи из map[]
func GetKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// получаем значения из map[]
func GetValues[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}
