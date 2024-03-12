package scanner

import "time"

const CheckpointFileExtension string = ".checkpoint"
const CheckpointRefreshInterval time.Duration = ScanRefreshInterval * 2

const IgnoreReasonDefault string = "ignored_by_default"
const IgnoreReasonDirPath string = "directory_path"
const IgnoreReasonFileExtensionIgnoredByConfig string = "file_extension_ignored_by_config"
const IgnoreReasonFileExtensionIgnoredByPolicy string = "file_extension_ignored_by_policy"
const IgnoreReasonFileExtensionNotIncluded string = "file_extension_not_included"
const IgnoreReasonFileIsBinary string = "file_is_binary"
const IgnoreReasonFileIsEmpty string = "file_is_empty"
const IgnoreReasonFileObjectPointerNil string = "file_object_pointer_nil"
const IgnoreReasonFileName string = "file_name"
const IgnoreReasonFilePath string = "file_path"

const ScanRefreshInterval time.Duration = time.Second * 5
