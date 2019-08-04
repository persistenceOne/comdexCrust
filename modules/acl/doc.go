package acl

// kv stores path is changed
// acl:{ZoneID}:{OrganizationID}:{Address} => ACL
// added validation in keeper of SetZoneAddress: check whether already exist of ZoneID
// added GetZones, GetOrganizations and GetOrganizationsByZoneID Query in keeper
// added organization struct : to string method and new organization create method
