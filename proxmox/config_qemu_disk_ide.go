package proxmox

import "strconv"

type QemuIdeDisk struct {
	AsyncIO    QemuDiskAsyncIO   `json:"asyncio,omitempty"`
	Backup     bool              `json:"backup"`
	Bandwidth  QemuDiskBandwidth `json:"bandwidth,omitempty"`
	Cache      QemuDiskCache     `json:"cache,omitempty"`
	Discard    bool              `json:"discard"`
	EmulateSSD bool              `json:"emulatessd"`
	Format     QemuDiskFormat    `json:"format"`
	Id         *uint             `json:"id,omitempty"`
	Replicate  bool              `json:"replicate"`
	Serial     QemuDiskSerial    `json:"serial,omitempty"`
	Size       uint              `json:"size"`
	Storage    string            `json:"storage"`
}

func (disk *QemuIdeDisk) convertDataStructure() *qemuDisk {
	return &qemuDisk{
		AsyncIO:    disk.AsyncIO,
		Backup:     disk.Backup,
		Bandwidth:  disk.Bandwidth,
		Cache:      disk.Cache,
		Discard:    disk.Discard,
		Disk:       true,
		EmulateSSD: disk.EmulateSSD,
		Format:     disk.Format,
		Id:         disk.Id,
		Replicate:  disk.Replicate,
		Serial:     disk.Serial,
		Size:       disk.Size,
		Storage:    disk.Storage,
		Type:       ide,
	}
}

func (disk QemuIdeDisk) Validate() error {
	return disk.convertDataStructure().validate()
}

type QemuIdeDisks struct {
	Disk_0 *QemuIdeStorage `json:"0,omitempty"`
	Disk_1 *QemuIdeStorage `json:"1,omitempty"`
	Disk_2 *QemuIdeStorage `json:"2,omitempty"`
	Disk_3 *QemuIdeStorage `json:"3,omitempty"`
}

// TODO write test
func (disks QemuIdeDisks) mapToApiValues(currentDisks *QemuIdeDisks, vmID uint, params map[string]interface{}, changes *qemuUpdateChanges) {
	tmpCurrentDisks := QemuIdeDisks{}
	if currentDisks != nil {
		tmpCurrentDisks = *currentDisks
	}
	diskMap := disks.mapToIntMap()
	currentDiskMap := tmpCurrentDisks.mapToIntMap()
	for i := range diskMap {
		diskMap[i].convertDataStructure().markDiskChanges(currentDiskMap[i].convertDataStructure(), vmID, "ide"+strconv.Itoa(int(i)), params, changes)
	}
}

func (disks QemuIdeDisks) mapToIntMap() map[uint8]*QemuIdeStorage {
	return map[uint8]*QemuIdeStorage{
		0: disks.Disk_0,
		1: disks.Disk_1,
		2: disks.Disk_2,
		3: disks.Disk_3,
	}
}

// TODO write test
func (QemuIdeDisks) mapToStruct(params map[string]interface{}) *QemuIdeDisks {
	disks := QemuIdeDisks{}
	var structPopulated bool
	if _, isSet := params["ide0"]; isSet {
		disks.Disk_0 = QemuIdeStorage{}.mapToStruct(params["ide0"].(string))
		structPopulated = true
	}
	if _, isSet := params["ide1"]; isSet {
		disks.Disk_1 = QemuIdeStorage{}.mapToStruct(params["ide1"].(string))
		structPopulated = true
	}
	if _, isSet := params["ide2"]; isSet {
		disks.Disk_2 = QemuIdeStorage{}.mapToStruct(params["ide2"].(string))
		structPopulated = true
	}
	if _, isSet := params["ide3"]; isSet {
		disks.Disk_3 = QemuIdeStorage{}.mapToStruct(params["ide3"].(string))
		structPopulated = true
	}
	if structPopulated {
		return &disks
	}
	return nil
}

func (disks QemuIdeDisks) Validate() error {
	diskMap := disks.mapToIntMap()
	for _, e := range diskMap {
		if e != nil {
			err := e.Validate()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type QemuIdePassthrough struct {
	AsyncIO    QemuDiskAsyncIO   `json:"asyncio,omitempty"`
	Backup     bool              `json:"backup"`
	Bandwidth  QemuDiskBandwidth `json:"bandwidth,omitempty"`
	Cache      QemuDiskCache     `json:"cache,omitempty"`
	Discard    bool              `json:"discard"`
	EmulateSSD bool              `json:"emulatessd"`
	File       string            `json:"file"`
	Replicate  bool              `json:"replicate"`
	Serial     QemuDiskSerial    `json:"serial,omitempty"`
	Size       uint              `json:"size"` //size is only returned and setting it has no effect
}

func (passthrough *QemuIdePassthrough) convertDataStructure() *qemuDisk {
	return &qemuDisk{
		AsyncIO:    passthrough.AsyncIO,
		Backup:     passthrough.Backup,
		Bandwidth:  passthrough.Bandwidth,
		Cache:      passthrough.Cache,
		Discard:    passthrough.Discard,
		EmulateSSD: passthrough.EmulateSSD,
		File:       passthrough.File,
		Replicate:  passthrough.Replicate,
		Serial:     passthrough.Serial,
		Type:       ide,
	}
}

func (passthrough QemuIdePassthrough) Validate() error {
	return passthrough.convertDataStructure().validate()
}

type QemuIdeStorage struct {
	CdRom       *QemuCdRom          `json:"cdrom,omitempty"`
	CloudInit   *QemuCloudInitDisk  `json:"cloudinit,omitempty"`
	Disk        *QemuIdeDisk        `json:"disk,omitempty"`
	Passthrough *QemuIdePassthrough `json:"passthrough,omitempty"`
}

// TODO write test
// converts to qemuStorage
func (storage *QemuIdeStorage) convertDataStructure() *qemuStorage {
	if storage == nil {
		return nil
	}
	generalizedStorage := qemuStorage{
		CdRom:     storage.CdRom,
		CloudInit: storage.CloudInit,
	}
	if storage.Disk != nil {
		generalizedStorage.Disk = storage.Disk.convertDataStructure()
	}
	if storage.Passthrough != nil {
		generalizedStorage.Passthrough = storage.Passthrough.convertDataStructure()
	}
	return &generalizedStorage
}

// TODO write test
func (QemuIdeStorage) mapToStruct(param string) *QemuIdeStorage {
	settings := splitStringOfSettings(param)
	tmpCdRom := qemuCdRom{}.mapToStruct(settings)
	if tmpCdRom != nil {
		if tmpCdRom.CdRom {
			return &QemuIdeStorage{CdRom: QemuCdRom{}.mapToStruct(*tmpCdRom)}
		} else {
			return &QemuIdeStorage{CloudInit: QemuCloudInitDisk{}.mapToStruct(*tmpCdRom)}
		}
	}

	tmpDisk := qemuDisk{}.mapToStruct(settings)
	if tmpDisk == nil {
		return nil
	}
	if tmpDisk.File == "" {
		return &QemuIdeStorage{Disk: &QemuIdeDisk{
			AsyncIO:    tmpDisk.AsyncIO,
			Backup:     tmpDisk.Backup,
			Bandwidth:  tmpDisk.Bandwidth,
			Cache:      tmpDisk.Cache,
			Discard:    tmpDisk.Discard,
			EmulateSSD: tmpDisk.EmulateSSD,
			Format:     tmpDisk.Format,
			Id:         tmpDisk.Id,
			Replicate:  tmpDisk.Replicate,
			Serial:     tmpDisk.Serial,
			Size:       tmpDisk.Size,
			Storage:    tmpDisk.Storage,
		}}
	}
	return &QemuIdeStorage{Passthrough: &QemuIdePassthrough{
		AsyncIO:    tmpDisk.AsyncIO,
		Backup:     tmpDisk.Backup,
		Bandwidth:  tmpDisk.Bandwidth,
		Cache:      tmpDisk.Cache,
		Discard:    tmpDisk.Discard,
		EmulateSSD: tmpDisk.EmulateSSD,
		File:       tmpDisk.File,
		Replicate:  tmpDisk.Replicate,
		Serial:     tmpDisk.Serial,
		Size:       tmpDisk.Size,
	}}
}

func (storage QemuIdeStorage) Validate() (err error) {
	// First check if more than one item is nil
	var subTypeSet bool
	if storage.CdRom != nil {
		subTypeSet = true
	}
	if storage.CloudInit != nil {
		if err = diskSubtypeSet(subTypeSet); err != nil {
			return
		}
		subTypeSet = true
	}
	if storage.Disk != nil {
		if err = diskSubtypeSet(subTypeSet); err != nil {
			return
		}
		subTypeSet = true
	}
	if storage.Passthrough != nil {
		if err = diskSubtypeSet(subTypeSet); err != nil {
			return
		}
	}
	// Validate sub items
	if storage.CdRom != nil {
		if err = storage.CdRom.Validate(); err != nil {
			return
		}
	}
	if storage.CloudInit != nil {
		if err = storage.CloudInit.Validate(); err != nil {
			return
		}
	}
	if storage.Disk != nil {
		if err = storage.Disk.Validate(); err != nil {
			return
		}
	}
	if storage.Passthrough != nil {
		err = storage.Passthrough.Validate()
	}
	return
}
