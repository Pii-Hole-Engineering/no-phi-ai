package rrr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	test_chunk_line_text_1 = "Lorem ipsum dolor sit amet"
	test_chunk_line_text_2 = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras viverra imperdiet orci, non laoreet nibh volutpat at. Donec ac dolor leo. Donec accumsan placerat leo dapibus tincidunt. Aenean egestas semper tellus nec dignissim. Maecenas blandit quam in lectus luctus egestas. Duis a mattis arcu. Vivamus interdum lacus nisi, a varius neque viverra ac. Nulla sed ligula nec tortor aliquam condimentum quis vitae lacus. In hac habitasse platea dictumst.\nCras dictum sapien turpis, condimentum tempus felis consequat ut. Donec blandit maximus dictum. Integer tincidunt malesuada odio, ut fringilla ligula pharetra nec. Morbi ac mattis metus. Nullam vehicula augue lorem, in eleifend nulla blandit quis. Donec pharetra dui a mattis malesuada. Fusce consectetur, ante at semper tempus, mauris lectus iaculis urna, porta hendrerit metus lectus a nibh. Etiam congue eros nec commodo mattis. Donec id ex ut quam laoreet mollis ut eget nunc. Pellentesque pulvinar hendrerit laoreet.\nVivamus leo metus, tempus semper varius et, porta blandit turpis. Aliquam faucibus dolor a velit consequat, et elementum ipsum aliquet. Aliquam semper, ligula consectetur blandit elementum, turpis orci hendrerit lorem, eu tincidunt risus eros nec tortor. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent elementum fringilla porta. Vivamus nisl justo, accumsan eget molestie sit amet, laoreet sit amet est. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Vestibulum mollis tellus sed venenatis fermentum. Integer vitae urna vitae massa bibendum semper eget at libero. Sed tempor et urna eu lacinia. Pellentesque magna ligula, molestie ac odio tempus, tincidunt malesuada metus. Etiam eu ex sodales, rhoncus dolor porttitor, dictum tellus. Curabitur condimentum lectus et feugiat faucibus. Vestibulum a suscipit ipsum.\nNullam interdum mi a sagittis feugiat. Sed nec consequat arcu. Nullam at pharetra tellus. Aliquam vitae justo feugiat, pellentesque enim id, pharetra purus. Interdum et malesuada fames ac ante ipsum primis in faucibus. Nulla consectetur velit nec molestie sagittis. In tincidunt commodo feugiat. Duis felis sem, semper a egestas in, commodo vitae felis. Curabitur blandit et enim eu vestibulum. Pellentesque imperdiet sodales magna sed porttitor. Praesent mollis placerat fermentum. Suspendisse sollicitudin dignissim tortor eget sollicitudin. Vivamus elementum erat leo, id mollis massa pulvinar placerat. Pellentesque vel augue id justo consequat pellentesque ut nec orci.\nNunc posuere nisi pellentesque, dignissim sem vel, rhoncus enim. Maecenas cursus dui id odio lobortis pulvinar. Phasellus pellentesque malesuada fermentum. Proin porttitor vel purus at accumsan. Aenean convallis mauris sed nisi interdum, nec consequat augue tincidunt. Ut at elit sit amet tellus lobortis rhoncus sodales ut nisi. Praesent finibus lorem vitae magna porttitor, ut finibus lectus commodo. Etiam in eros felis. Nullam molestie ultrices magna nec suscipit. Fusce vitae pharetra ante. Cras iaculis magna nisi, ut accumsan ex sagittis quis. Integer mauris elit, pretium non volutpat vel, scelerisque at lectus. In hac habitasse platea dictumst. Quisque pretium libero id nisi rutrum, et auctor velit maximus.\nPellentesque ut nisl ex. Pellentesque vitae purus nibh. Etiam ornare magna ut ipsum ullamcorper tristique vitae ornare mauris. Nam eleifend metus nec risus porttitor eleifend. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed venenatis, magna eget vulputate pretium, tellus libero viverra nunc, vitae commodo neque urna quis lorem. Nullam nec sem a risus malesuada ullamcorper. Quisque viverra eros ut libero consequat sollicitudin. Praesent fringilla nisi vel metus auctor commodo. Ut lacinia rutrum magna, sit amet lacinia ante pulvinar in.\nNam in tristique mauris. Integer lectus felis, vehicula a pretium sodales, lobortis ac mi. Nam et pulvinar massa. Quisque lectus tortor, mattis a ipsum et, porta efficitur purus. Integer rhoncus velit nisl, eu tristique metus volutpat condimentum. In mi est, tincidunt quis tincidunt id, sagittis eget justo. Maecenas erat diam, placerat ut rhoncus quis, semper id dolor. Nullam varius luctus ligula. Aenean vel ullamcorper velit. Interdum et malesuada fames ac ante ipsum primis in faucibus. Cras nulla est, consequat eget tempus sed, luctus id nulla. Pellentesque a erat tellus. Suspendisse potenti. Ut feugiat risus nec congue sodales. Cras at porttitor orci, ut molestie enim.\nFusce nisl metus, molestie sed eleifend sit amet, interdum nec massa. Nunc ac purus lobortis, malesuada diam maximus, porttitor quam. Duis non consectetur ante, viverra eleifend est. Quisque sit amet eros a nisi accumsan maximus ut eu orci. Sed tellus neque, suscipit eget justo a, pharetra ultrices urna. Nam mollis erat vel quam faucibus, vitae cursus justo dignissim. Mauris sodales est vitae vehicula laoreet. In ornare diam sit amet purus mollis, quis tincidunt dolor bibendum.\nAliquam ultricies purus ac aliquet posuere. Sed efficitur ultricies magna, eget imperdiet leo dictum non. Ut vestibulum nulla quis sodales porttitor. Aliquam tincidunt purus sit amet mattis accumsan. Proin id auctor odio. Nam vitae nunc commodo, elementum risus sit amet, ultricies est. Donec odio magna, convallis at venenatis et, varius nec orci. Suspendisse sapien mi, placerat nec ullamcorper ut, consectetur et erat. Fusce dapibus condimentum nunc et egestas. Proin malesuada posuere malesuada. Duis rutrum elit quis ultrices tempus. Vestibulum eros sem, auctor a lorem sed, vulputate ultrices quam.\nPhasellus id semper augue. Etiam eleifend lacus quam, et tempus orci efficitur eget. Sed varius pellentesque porttitor. Aenean et luctus enim. Suspendisse eu ipsum sollicitudin, vestibulum lacus ut, venenatis sapien. Fusce scelerisque erat nibh, in consequat nisi aliquet a. Nullam tellus mi, consequat vel malesuada sed, placerat at nunc. Nulla facilisi. Duis pellentesque lacinia lorem, non rutrum justo sollicitudin dapibus. Proin molestie dui pharetra ante viverra auctor. Suspendisse vel erat in arcu scelerisque eleifend eu blandit nisl.\nVivamus et mauris accumsan, suscipit massa sed, consectetur ipsum. Sed lobortis neque tristique metus vulputate, nec mattis mi fermentum. Maecenas a elit quis ipsum convallis ultricies eget sit amet lacus. In interdum dictum lacus ut venenatis. Aenean cursus, dui ut vulputate dapibus, quam neque interdum felis, in tristique nibh felis eu elit. Nunc et varius justo, a interdum risus. Donec ornare vel elit non efficitur. Suspendisse vel nisi gravida, pretium elit vitae, laoreet urna. Nam consequat tincidunt erat, in placerat nibh condimentum nec. Donec vel auctor dui. Cras id magna non ante finibus posuere ac vitae metus. Suspendisse lobortis, nisi sed convallis tempus, purus."
)

// TestChunkFileToRequests unit test function tests the
// ChunkFileToRequests() function.
func TestChunkFileToRequests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expected_requests []Request
		expected_error    error
		in                ChunkFileInput
		name              string
	}{
		{
			expected_requests: []Request{
				{
					MetadataRequestResponse: MetadataRequestResponse{
						Commit: MetadataRequestResponseCommit{
							ID: "test_commit",
						},
						Object: MetadataRequestResponseObject{
							ID: "test_object",
						},
						Repository: MetadataRequestResponseRepository{
							ID: "test_repo",
						},
					},
					Text: "Lorem ipsum dolor sit amet",
				},
			},
			expected_error: ErrChunkFileToRequestsInFileNil,
			in: ChunkFileInput{
				RepoID:       "test_repo",
				CommitID:     "test_commit",
				File:         nil, // TODO
				MaxChunkSize: 100,
			},
			name: "Pass_1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requests, err := ChunkFileToRequests(tt.in)
			if tt.expected_error != nil {
				assert.ErrorIs(t, err, tt.expected_error)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, tt.expected_requests, requests)
		})
	}
}

// TestChunkLineToRequests unit test function tests the
// ChunkLineToRequests() function.
func TestChunkLineToRequests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expected_error        error
		expected_final_offset int
		expected_num_requests int
		in                    ChunkLineInput
		name                  string
	}{
		{
			expected_error:        nil,
			expected_final_offset: len(test_chunk_line_text_1),
			expected_num_requests: 4,
			in: ChunkLineInput{
				CommitID:     "test_commit",
				Line:         test_chunk_line_text_1,
				MaxChunkSize: 10,
				Offset:       0,
				RepoID:       "test_repo",
				ObjectID:     "test_object",
			},
			name: "Pass_1",
		},
		{
			expected_error:        nil,
			expected_final_offset: len(test_chunk_line_text_1),
			expected_num_requests: 3,
			in: ChunkLineInput{
				CommitID:     "test_commit",
				Line:         test_chunk_line_text_1,
				MaxChunkSize: 12,
				Offset:       0,
				RepoID:       "test_repo",
				ObjectID:     "test_object",
			},
			name: "Pass_2",
		},
		{
			expected_error:        nil,
			expected_final_offset: len(test_chunk_line_text_1),
			expected_num_requests: 1,
			in: ChunkLineInput{
				CommitID:     "test_commit",
				Line:         test_chunk_line_text_1,
				MaxChunkSize: 30,
				Offset:       0,
				RepoID:       "test_repo",
				ObjectID:     "test_object",
			},
			name: "Pass_3",
		},
		{
			expected_error:        nil,
			expected_final_offset: 100,
			expected_num_requests: 0,
			in: ChunkLineInput{
				CommitID:     "test_commit",
				Line:         "",
				MaxChunkSize: 50,
				Offset:       100,
				RepoID:       "test_repo",
				ObjectID:     "test_object",
			},
			name: "Pass_4",
		},
		{
			expected_error:        nil,
			expected_final_offset: len(test_chunk_line_text_2),
			expected_num_requests: 7,
			in: ChunkLineInput{
				CommitID:     "test_commit",
				Line:         test_chunk_line_text_2,
				MaxChunkSize: 1000,
				Offset:       0,
				RepoID:       "test_repo",
				ObjectID:     "test_object",
			},
			name: "Pass_5",
		},
		{
			expected_error:        nil,
			expected_final_offset: len(test_chunk_line_text_2) + 333,
			expected_num_requests: 7,
			in: ChunkLineInput{
				CommitID:     "test_commit",
				Line:         test_chunk_line_text_2,
				MaxChunkSize: 1000,
				Offset:       333,
				RepoID:       "test_repo",
				ObjectID:     "test_object",
			},
			name: "Pass_6",
		},
		{
			expected_error:        ErrNewRequestEmptyCommitID,
			expected_final_offset: 100,
			expected_num_requests: 0,
			in: ChunkLineInput{
				CommitID:     "",
				Line:         test_chunk_line_text_1,
				MaxChunkSize: 50,
				Offset:       100,
				RepoID:       "test_repo",
				ObjectID:     "test_object",
			},
			name: "Err_CommitID",
		},
		{
			expected_error:        ErrNewRequestEmptyObjectID,
			expected_final_offset: 100,
			expected_num_requests: 0,
			in: ChunkLineInput{
				CommitID:     "test_commit",
				Line:         test_chunk_line_text_1,
				MaxChunkSize: 50,
				Offset:       100,
				RepoID:       "test_repo",
				ObjectID:     "",
			},
			name: "Err_ObjectID",
		},
		{
			expected_error:        ErrNewRequestEmptyRepositoryID,
			expected_final_offset: 100,
			expected_num_requests: 0,
			in: ChunkLineInput{
				CommitID:     "test_commit",
				Line:         test_chunk_line_text_1,
				MaxChunkSize: 50,
				Offset:       100,
				RepoID:       "",
				ObjectID:     "test_object",
			},
			name: "Err_RepoID",
		},
		{
			expected_error:        ErrMaxChunkSizeInvalid,
			expected_final_offset: 0,
			expected_num_requests: 0,
			in: ChunkLineInput{
				CommitID:     "test_commit",
				Line:         test_chunk_line_text_1,
				MaxChunkSize: -1,
				Offset:       0,
				RepoID:       "test_repo",
				ObjectID:     "test_object",
			},
			name: "Err_ChunkSize_Negative",
		},
		{
			expected_error:        ErrMaxChunkSizeInvalid,
			expected_final_offset: 0,
			expected_num_requests: 0,
			in: ChunkLineInput{
				CommitID:     "test_commit",
				Line:         test_chunk_line_text_1,
				MaxChunkSize: 0,
				Offset:       0,
				RepoID:       "test_repo",
				ObjectID:     "test_object",
			},
			name: "Err_ChunkSize_Zero",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			final_offset, requests, err := ChunkLineToRequests(test.in)
			if test.expected_error == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, test.expected_error, err)
			}

			for _, request := range requests {
				assert.Equal(t, test.in.RepoID, request.Repository.ID)
			}
			assert.Equal(t, test.expected_num_requests, len(requests))
			assert.Equal(t, test.expected_final_offset, final_offset)
		})
	}

	// TODO
}
