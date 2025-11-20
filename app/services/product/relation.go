package service_product

func (repo *Service) SyncRelations(mainID uint, ids []uint) error {
	return repo.repositories.product.SyncRelations(mainID, ids)
}
